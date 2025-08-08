package tools

type ModelRegistryTools struct {
	Name        string
	Description string
	Handler     any
}

var Tools = []ModelRegistryTools{
	{
		Name:        "register_model",
		Description: "Register a model in the model registry",
		Handler:     RegisterModelHandler,
	},
	{
		Name:        "create_model_version",
		Description: "Create a model version related to a registered model in the model registry",
		Handler:     CreateModelVersionHandler,
	},
	{
		Name:        "create_model_artifact",
		Description: "Create a model artifact related to a model version in the model registry",
		Handler:     CreateModelArtifactHandler,
	},
}
