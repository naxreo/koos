package main

import (
	"fmt"

	"github.com/naxreo/koos/kube"
	"github.com/naxreo/koos/rha"
	"k8s.io/client-go/kubernetes"
)

const (
	initPasswd  = "ChangePassword"
	redisPasswd = "RHAPassword"
)

type RHAInst struct {
	masterName         string
	slaveName          string
	namespace          string
	masterHeadlessName string
	slaveHeadlessName  string
	svcName            string
}

func deployKube(clientset *kubernetes.Clientset, inst *RHAInst) bool {
	isInit := true

	// for master
	if !kube.IsDeployment(clientset, inst.namespace, inst.masterName) {
		ok := rha.CreateRHADeployment(clientset, inst.namespace, inst.masterName, "")
		if !ok {
			panic("Failed to create for the master deployment")
		}
		fmt.Printf("Created %q...\n", inst.masterName)
	} else {
		isInit = false
	}

	if !kube.IsService(clientset, inst.namespace, inst.masterHeadlessName) {
		ok := rha.CreateRHAHeadless(clientset, inst.namespace, inst.masterHeadlessName, inst.masterName)
		if !ok {
			panic("Failed to create for the master headless")
		}
		fmt.Printf("Created %q...\n", inst.masterHeadlessName)
	} else {
		isInit = false
	}

	// for worker
	if !kube.IsDeployment(clientset, inst.namespace, inst.slaveName) {
		ok := rha.CreateRHADeployment(clientset, inst.namespace, inst.slaveName, "")
		if !ok {
			panic("Failed to create for the worker deployment")
		}
		fmt.Printf("Created %q...\n", inst.slaveName)
	} else {
		isInit = false
	}
	if !kube.IsService(clientset, inst.namespace, inst.slaveHeadlessName) {
		ok := rha.CreateRHAHeadless(clientset, inst.namespace, inst.slaveHeadlessName, inst.slaveName)
		if !ok {
			panic("Failed to create for the worker headless")
		}
		fmt.Printf("Created %q...\n", inst.slaveHeadlessName)
	} else {
		isInit = false
	}

	// ClusterIP Service for initialization
	if !kube.IsService(clientset, inst.namespace, inst.svcName) {
		ok := rha.CreateRHAService(clientset, inst.namespace, inst.svcName, inst.masterName)
		if !ok {
			panic("Failed to create for the service")
		}
		fmt.Printf("Created %q...\n", inst.svcName)
	} else {
		isInit = false
	}

	if isInit {
		fmt.Println("End of RHA initialization...")
	} else {
		fmt.Println("Restart controller...")
	}
	return isInit
}

func main() {

	fmt.Println("Start initialization...")
	clientset := kube.Config()
	inst := &RHAInst{"rha-1", "rha-2", "default", "rha-1-headless", "rha-2-headless", "rha-svc"}

	maddr := inst.masterHeadlessName + "." + inst.namespace + ".svc.cluster.local:6379"
	// saddr := inst.slaveHeadlessName + "." + inst.namespace + ".svc.cluster.local:6379"
	maddr = "localhost:6379"

	//	Init RHA
	var m *rha.RedisPod

	if deployKube(clientset, inst) {
		// redis master & slave code here
		fmt.Println("Configure a master and slave...")
	}

	m = &rha.RedisPod{"master", true, maddr, initPasswd, 0}
	mr := rha.Connect(m)
	// s := &rha.RedisPod{"master", true, saddr, initPasswd, 0}
	// sr := rha.Connect(m)

	// ctx := context.Background()
	// for i := 0; i < 10; i++ {
	// 	mping, err := mr.Ping(ctx).Result()
	// 	if err != nil {
	// 		fmt.Printf("master redis connect failed %s\n", err.Error())
	// 	}
	// 	fmt.Printf("%q\n", mping)

	// 	info := mr.Info(ctx, "")
	// 	fmt.Printf("%q\n", info)
	// 	role := mr.ConfigGet(ctx, "role")
	// 	fmt.Printf("%q\n", role)

	// 	time.Sleep(1 * time.Second)
	// }

	p := rha.RhaPing(mr, 3)
	fmt.Println(p)
}
