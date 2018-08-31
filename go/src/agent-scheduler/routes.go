package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.hpls.local/ppsc/agent-scheduler/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"POST",
		"/submitjob",
		handlers.SubmitJob,
	},
	Route{
		"Index",
		"GET",
		"/testjob",
		handlers.TestJob,
	},

	Route{
		"Index",
		"POST",
		"/api/v1/job/request/submit",
		handlers.CreateAndSubmitJobScript,
	},
	Route{
		"qstat-Q",
		"POST",
		"/queryQueueLimits",
		handlers.QueryQueueLimit,
	},
	Route{
		"qstat-B",
		"POST",
		"/serverSummary",
		handlers.ServerSummary,
	},
	Route{
		"qstat-f",
		"POST",
		"/JobQueryByJobId",
		handlers.JobQueryByJobId,
	},
	Route{
		"qstat-a",
		"POST",
		"/GetJobInfo",
		handlers.GetJobInfo,
	},
	Route{
		"qstat",
		"POST",
		"/GetQueueNJobInfo",
		handlers.GetQueueNJobInfo,
	},
	Route{
		"qstat-iu-",
		"POST",
		"/GetQueuedJobInfoByUserId",
		handlers.GetQueuedJobInfoByUserId,
	},
	Route{
		"qstat-u-",
		"POST",
		"/GetJobInfoByUserId",
		handlers.GetJobInfoByUserId,
	},
	Route{
		"qstat-s-1",
		"POST",
		"/GetJobInfoWithComment",
		handlers.GetJobInfoWithComment,
	},
	Route{
		"qstat-n-1-",
		"POST",
		"/GetNodeInfoByJobId",
		handlers.GetNodeInfoByJobId,
	},
	Route{
		"pbsnodes-a",
		"POST",
		"/GetAllPbsNodesInfo",
		handlers.GetAllPbsNodesInfo,
	},
	Route{
		"pbsnodes-a",
		"POST",
		"/GetClusterInfo",
		handlers.GetClusterInfo,
	},
	Route{
		"qstat-Qf",
		"POST",
		"/QueueQueryByQueueId",
		handlers.QueueQueryByQueueId,
	},
	Route{
		"Index",
		"POST",
		"/api/v1/job/request/terminate",
		handlers.TerminateJobHandler,
	},
}
