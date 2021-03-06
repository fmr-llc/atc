package gc

import (
	"errors"
	"time"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/metric"
	"github.com/concourse/atc/worker"
)

const HijackedContainerTimeout = 5 * time.Minute

var containerCollectorFailedErr = errors.New("container collector failed")

type containerCollector struct {
	rootLogger          lager.Logger
	containerRepository db.ContainerRepository
	jobRunner           WorkerJobRunner
}

func NewContainerCollector(
	logger lager.Logger,
	containerRepository db.ContainerRepository,
	jobRunner WorkerJobRunner,
) Collector {
	return &containerCollector{
		rootLogger:          logger,
		containerRepository: containerRepository,
		jobRunner:           jobRunner,
	}
}

type job struct {
	JobName string
	RunFunc func(worker.Worker)
}

func (j *job) Name() string {
	return j.JobName
}

func (j *job) Run(w worker.Worker) {
	j.RunFunc(w)
}

func (c *containerCollector) Run() error {
	logger := c.rootLogger.Session("run")

	logger.Debug("start")
	defer logger.Debug("done")

	var err error

	orphanedErr := c.cleanupOrphanedContainers(logger.Session("orphaned-containers"))
	if orphanedErr != nil {
		c.rootLogger.Error("container-collector", orphanedErr)
		err = containerCollectorFailedErr
	}

	failedErr := c.cleanupFailedContainers(logger.Session("failed-containers"))
	if failedErr != nil {
		c.rootLogger.Error("container-collector", failedErr)
		err = containerCollectorFailedErr
	}

	return err
}

func (c *containerCollector) cleanupFailedContainers(logger lager.Logger) error {

	failedContainers, err := c.containerRepository.FindFailedContainers()
	if err != nil {
		logger.Error("failed-to-find-failed-containers-for-deletion", err)
		return err
	}

	failedContainerHandles := []string{}
	var failedContainerstoDestroy = []destroyableContainer{}

	if len(failedContainers) > 0 {
		for _, container := range failedContainers {
			failedContainerHandles = append(failedContainerHandles, container.Handle())
			failedContainerstoDestroy = append(failedContainerstoDestroy, container)
		}
	}

	logger.Debug("found-failed-containers-for-deletion", lager.Data{
		"failed-containers": failedContainerHandles,
	})

	metric.FailedContainersToBeGarbageCollected{
		Containers: len(failedContainerHandles),
	}.Emit(logger)

	destroyDBContainers(logger, failedContainerstoDestroy)

	return nil
}

func (c *containerCollector) cleanupOrphanedContainers(logger lager.Logger) error {
	creatingContainers, createdContainers, destroyingContainers, err := c.containerRepository.FindOrphanedContainers()

	if err != nil {
		logger.Error("failed-to-get-orphaned-containers-for-deletion", err)
		return err
	}

	creatingContainerHandles := []string{}
	createdContainerHandles := []string{}
	destroyingContainerHandles := []string{}

	if len(creatingContainers) > 0 {
		for _, container := range creatingContainers {
			creatingContainerHandles = append(creatingContainerHandles, container.Handle())
		}
	}

	if len(createdContainers) > 0 {
		for _, container := range createdContainers {
			createdContainerHandles = append(createdContainerHandles, container.Handle())
		}
	}

	if len(destroyingContainers) > 0 {
		for _, container := range destroyingContainers {
			destroyingContainerHandles = append(destroyingContainerHandles, container.Handle())
		}
	}

	logger.Debug("found-orphaned-containers-for-deletion", lager.Data{
		"creating-containers":   creatingContainerHandles,
		"created-containers":    createdContainerHandles,
		"destroying-containers": destroyingContainerHandles,
	})

	metric.CreatingContainersToBeGarbageCollected{
		Containers: len(creatingContainerHandles),
	}.Emit(logger)

	metric.CreatedContainersToBeGarbageCollected{
		Containers: len(createdContainerHandles),
	}.Emit(logger)

	metric.DestroyingContainersToBeGarbageCollected{
		Containers: len(destroyingContainerHandles),
	}.Emit(logger)

	var workerCreatedContainers = make(map[string][]db.CreatedContainer)

	for _, createdContainer := range createdContainers {
		containers, ok := workerCreatedContainers[createdContainer.WorkerName()]
		if ok {
			// update existing array
			containers = append(containers, createdContainer)
			workerCreatedContainers[createdContainer.WorkerName()] = containers
		} else {
			// create new array
			workerCreatedContainers[createdContainer.WorkerName()] = []db.CreatedContainer{createdContainer}
		}
	}

	for worker, createdContainers := range workerCreatedContainers {
		// prevent closure from capturing last value of loop
		c.jobRunner.Try(logger,
			worker,
			&job{
				JobName: worker,
				RunFunc: destroyCreatedContainers(logger, createdContainers),
			},
		)
	}

	var workerContainers = make(map[string][]db.DestroyingContainer)

	for _, destroyingContainer := range destroyingContainers {
		containers, ok := workerContainers[destroyingContainer.WorkerName()]
		if ok {
			// update existing array
			containers = append(containers, destroyingContainer)
			workerContainers[destroyingContainer.WorkerName()] = containers
		} else {
			// create new array
			workerContainers[destroyingContainer.WorkerName()] = []db.DestroyingContainer{destroyingContainer}
		}
	}

	for worker, destroyingContainers := range workerContainers {
		c.jobRunner.Try(logger,
			worker,
			&job{
				JobName: worker,
				RunFunc: destroyDestroyingContainers(logger, destroyingContainers),
			},
		)
	}
	return nil
}

func destroyCreatedContainers(logger lager.Logger, containers []db.CreatedContainer) func(worker.Worker) {
	return func(workerClient worker.Worker) {
		destroyingContainers := []db.DestroyingContainer{}

		for _, container := range containers {

			var destroyingContainer db.DestroyingContainer
			if container.IsHijacked() {
				cLog := logger.Session("mark-hijacked-container", lager.Data{
					"container": container.Handle(),
					"worker":    workerClient.Name(),
				})

				var err error
				destroyingContainer, err = markHijackedContainerAsDestroying(cLog, container, workerClient.GardenClient())
				if err != nil {
					cLog.Error("failed-to-transition", err)
					return
				}
			} else {
				cLog := logger.Session("mark-created-as-destroying", lager.Data{
					"container": container.Handle(),
					"worker":    workerClient.Name(),
				})
				var err error
				destroyingContainer, err = container.Destroying()
				if err != nil {
					cLog.Error("failed-to-transition", err)
					return
				}
			}
			if destroyingContainer != nil {
				destroyingContainers = append(destroyingContainers, destroyingContainer)
			}

		}
		tryToDestroyContainers(logger.Session("destroy-containers"), destroyingContainers, workerClient)
	}
}

func destroyDestroyingContainers(logger lager.Logger, containers []db.DestroyingContainer) func(worker.Worker) {
	return func(workerClient worker.Worker) {
		cLog := logger.Session("destroy-containers-on-worker", lager.Data{
			"worker": workerClient.Name(),
		})
		tryToDestroyContainers(cLog, containers, workerClient)
	}
}

func markHijackedContainerAsDestroying(
	logger lager.Logger,
	hijackedContainer db.CreatedContainer,
	gardenClient garden.Client,
) (db.DestroyingContainer, error) {

	gardenContainer, found, err := findContainer(gardenClient, hijackedContainer.Handle())
	if err != nil {
		logger.Error("failed-to-lookup-container-in-garden", err)
		return nil, err
	}

	if !found {
		logger.Debug("hijacked-container-not-found-in-garden")

		destroyingContainer, err := hijackedContainer.Destroying()
		if err != nil {
			logger.Error("failed-to-mark-container-as-destroying", err)
			return nil, err
		}
		return destroyingContainer, nil
	}

	err = gardenContainer.SetGraceTime(HijackedContainerTimeout)
	if err != nil {
		logger.Error("failed-to-set-grace-time-on-hijacked-container", err)
		return nil, err
	}

	_, err = hijackedContainer.Discontinue()
	if err != nil {
		logger.Error("failed-to-mark-container-as-destroying", err)
		return nil, err
	}

	return nil, nil
}

func tryToDestroyContainers(
	logger lager.Logger,
	containers []db.DestroyingContainer,
	workerClient worker.Worker,
) {
	logger.Debug("start")
	defer logger.Debug("done")

	gardenDeleteHandles := []string{}
	gardenDeleteContainers := []destroyableContainer{}

	dbDeleteContainers := []destroyableContainer{}

	gardenClient := workerClient.GardenClient()
	reaperClient := workerClient.ReaperClient()

	for _, container := range containers {
		if container.IsDiscontinued() {
			cLog := logger.Session("discontinued", lager.Data{"handle": container.Handle()})

			_, found, err := findContainer(gardenClient, container.Handle())
			if err != nil {
				cLog.Error("failed-to-lookup-container-in-garden", err)
			}

			if found {
				cLog.Debug("still-present-in-garden")
			} else {
				cLog.Debug("container-no-longer-present-in-garden")
				dbDeleteContainers = append(dbDeleteContainers, container)
			}
		} else {
			gardenDeleteHandles = append(gardenDeleteHandles, container.Handle())
			gardenDeleteContainers = append(gardenDeleteContainers, container)
		}
	}

	if len(gardenDeleteHandles) > 0 {
		err := reaperClient.DestroyContainers(gardenDeleteHandles)
		if err != nil {
			logger.Error("failed-to-destroy-garden-containers", err, lager.Data{"handlers": gardenDeleteHandles})
		} else {
			logger.Debug("completed-destroyed-in-garden", lager.Data{"handlers": gardenDeleteHandles})
			dbDeleteContainers = append(dbDeleteContainers, gardenDeleteContainers...)
		}
	}

	destroyDBContainers(logger, dbDeleteContainers)
	logger.Debug("destroyed-in-db")
}

type destroyableContainer interface {
	Destroy() (bool, error)
}

func destroyDBContainers(logger lager.Logger, dbContainers []destroyableContainer) {
	logger.Debug("destroying-start-in-db", lager.Data{"length": len(dbContainers)})
	defer logger.Debug("destroying-done-in-db")

	for _, dbContainer := range dbContainers {
		destroyed, err := dbContainer.Destroy()
		if err != nil {
			logger.Error("failed-to-destroy-database-container", err)
			continue
		}

		if !destroyed {
			logger.Info("could-not-destroy-database-container")
			continue
		}

		metric.ContainersDeleted.Inc()
		logger.Debug("destroyed-container-in-db")
	}
}

func findContainer(gardenClient garden.Client, handle string) (garden.Container, bool, error) {
	gardenContainer, err := gardenClient.Lookup(handle)
	if err != nil {
		if _, ok := err.(garden.ContainerNotFoundError); ok {
			return nil, false, nil
		}
		return nil, false, err
	}
	return gardenContainer, true, nil
}
