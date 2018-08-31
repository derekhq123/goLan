package handlers

type PBSScript struct {
	Id            int      `json:"id"`
	MetaID        string   `json:"metaID"`
	Name          string   `json:"name"`
	Resources     string   `json:"resources"`
	Chunk         string   `json:"chunk"`
	NCPUs         string   `json:"ncpu"`
	Mem           string   `json:"mem"`
	MPIProcs      string   `json:"mpiprocs"`
	OMPThreads    string   `json:"ompthreads"`
	WallTime      string   `json:"walltime"`
	OutputFile    string   `json:"outputfile"`
	ErrorFile     string   `json:"errorfile"`
	QueueName     string   `json:"queuename"`
	Description   string   `json:"description"`
	WorkDirectory string   `json:"workdirectory"`
	PreProcess    string   `json:"preprocess"`
	PostProcess   string   `json:"postprocess"`
	CustomScript  string   `json:"customscript"`
	Module        []string `json:"module"`
}
