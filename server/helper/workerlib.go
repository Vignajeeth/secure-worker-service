package helper

import (
	"bytes"
	"log"
	"os/exec"
	"time"
)

// CreateJob is a factory method used to create a job object when the command, user and jobid is  given.
func CreateJob(commandString, issuedUser string, jobid int) *Job {
	job := &Job{
		ID:        jobid, //#CHANGE
		Status:    StatusCode(Created),
		Command:   commandString,
		UserId:    issuedUser,
		StartTime: time.Now(),
		// StopTime:  nil,
	}
	JobDB = append(JobDB, *job)
	return job
}

// Workerlib is the exported function when using StartJob. It returns a job object immediately so that users
// can query or stop the job. It sends the job to a wrapper function which performs update operations
// in addition to running the job.
func Workerlib(jobBody Job, jobid int, user string) Job {

	j := CreateJob(jobBody.Command, user, jobid)

	go wrapperFunc(*j)

	return *j
}

// wrapperFunc performs status update operations on the DB and acts as a wrapper to the job process.
func wrapperFunc(j Job) {
	for i := range JobDB {
		if JobDB[i].ID == j.ID {
			JobDB[i].Status = StatusCode(Running)
			break
		}
	}

	output, err := LinuxJob(j.Command, j.ID)

	for i := range JobDB {
		if JobDB[i].ID == j.ID {
			if err != "" {
				JobDB[i].Status = StatusCode(Failed)
				JobDB[i].Result = err
				JobDB[i].StopTime = time.Now()
			} else {
				JobDB[i].Status = StatusCode(Completed)
				JobDB[i].Result = output
				JobDB[i].StopTime = time.Now()
			}

		}
	}
}

// LinuxJob is the base method which performs the job and handles exit.
func LinuxJob(commandString string, commandId int) (string, string) {
	cmd := exec.Command("bash", "-c", commandString)
	var cmdOutput, cmdError bytes.Buffer
	cmd.Stdout = &cmdOutput // Output of the process.
	cmd.Stderr = &cmdError  // Error of the process.

	err := cmd.Start()
	if err != nil {
		log.Println(err)
		log.Println("error in start")
		return cmdOutput.String(), cmdError.String()
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {

	case <-time.After(30 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Println("Error: Timeout reached, but unable to kill process of Job", commandId)
		} else {
			log.Println("Process killed as Timeout is reached for Job", commandId)
		}

	case <-Ch[commandId]:
		if err := cmd.Process.Kill(); err != nil {
			log.Println("Error: Stop job request received, but failed to kill process of Job", commandId)
		} else {
			log.Println("Process killed as it was manually stopped reached for Job", commandId)
		}

	case err := <-done:
		if err != nil {
			log.Println("Error: Process finished for Job with error.", commandId)
		} else {
			log.Println("Process finished successfully for Job", commandId)
		}
	}
	log.Println("\n", cmdOutput.String())
	// log.Println("err:", cmdError.String())
	return cmdOutput.String(), cmdError.String()

}

// KillJob stops the job by using a global channel indexed by the jobid.
func KillJob(jobid int) {
	for i := range JobDB {
		if JobDB[i].ID == jobid {
			go func() { Ch[jobid] <- true }()
			JobDB[i].StopTime = time.Now()
			break
		}
	}

}
