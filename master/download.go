package master

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/h2oai/steamY/lib/fs"
	"github.com/h2oai/steamY/master/az"
	"github.com/h2oai/steamY/srv/compiler"
	srvweb "github.com/h2oai/steamY/srv/web"
)

const (
	paramType      = "type"
	paramTypeModel = "model"
	paramArtifact  = "artifact"
	paramProjectId = "project-id"
	paramLabelName = "label-name"
	paramModelId   = "model-id"

	// model artifact types
	javaClass    = "java-class"     // foo.java
	javaClassDep = "java-class-dep" // gen-model.jar
	javaJar      = "java-jar"       // foo.jar
	javaWar      = "java-war"       // foo.war
)

type DownloadHandler struct {
	az                     az.Az
	workingDirectory       string
	webService             srvweb.Service
	compilerServiceAddress string
}

func newDownloadHandler(az az.Az, workingDirectory string, webService srvweb.Service, compilerServiceAddress string) *DownloadHandler {
	return &DownloadHandler{
		az,
		workingDirectory,
		webService,
		compilerServiceAddress,
	}
}

func (s *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("File download request received.")

	pz, azerr := s.az.Identify(r)
	if azerr != nil {
		log.Println(azerr)
		http.Error(w, fmt.Sprintf("Authentication failed: %s", azerr), http.StatusForbidden)
	}

	values := r.URL.Query()

	typ := values.Get(paramType)

	if len(typ) == 0 {
		http.Error(w, fmt.Sprintf("Missing %s", paramType), http.StatusBadRequest)
		return
	}

	switch typ {
	case paramTypeModel:
		artifact := values.Get(paramArtifact)

		if len(artifact) == 0 {
			http.Error(w, fmt.Sprintf("Missing %s", paramArtifact), http.StatusBadRequest)
			return
		}

		modelIdValue := values.Get(paramModelId)
		if len(modelIdValue) != 0 {
			modelId, err := strconv.ParseInt(modelIdValue, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Not a serial number %s=%s: %s", paramModelId, modelIdValue, err), http.StatusBadRequest)
				return
			}

			s.serveModel(w, r, pz, modelId, artifact)
			return
		}

		labelName := values.Get(paramLabelName)
		if len(labelName) != 0 {
			projectIdValue := values.Get(paramProjectId)
			if len(projectIdValue) == 0 {
				http.Error(w, fmt.Sprintf("Missing %s", paramProjectId), http.StatusBadRequest)
				return
			}

			projectId, err := strconv.ParseInt(projectIdValue, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Not a serial number %s=%s: %s", paramProjectId, projectIdValue, err), http.StatusBadRequest)
				return
			}

			if projectId <= 0 {
				http.Error(w, fmt.Sprintf("Invalid %s: %s", paramProjectId, projectId), http.StatusBadRequest)
				return
			}

			labels, err := s.webService.GetLabelsForProject(pz, projectId)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed reading labels for project %d: %s", projectId, err), http.StatusInternalServerError)
				return
			}

			for _, label := range labels {
				if label.Name == labelName {
					if label.ModelId > 0 {
						s.serveModel(w, r, pz, label.ModelId, artifact)
						return
					}
					http.Error(w, fmt.Sprintf("No model associated with label: %s", labelName), http.StatusNotFound)
					return
				}
			}

			http.Error(w, fmt.Sprintf("No label found: %s", labelName), http.StatusNotFound)
			return

		}
	default:
		http.Error(w, fmt.Sprintf("Invalid %s: %s", paramType, typ), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (s *DownloadHandler) serveModel(w http.ResponseWriter, r *http.Request, pz az.Principal, modelId int64, artifact string) {
	switch artifact {
	case javaClass, javaClassDep, javaJar, javaWar:
		// Call the API to get the model details.
		// We assume that if the GetModel() call succeeds, the principal has
		//   permissions and privileges to read this model, and consequently
		//   allowed to download it.
		model, err := s.webService.GetModel(pz, modelId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed reading model %d: %s", modelId, err), http.StatusUnauthorized)
			return
		}
		modelLocation := model.Location
		if len(modelLocation) == 0 {
			http.Error(w, fmt.Sprintf("Failed reading model %d: the model was not saved correctly", model.Id), http.StatusNotFound)
			return
		}

		var filePath string
		switch artifact {
		case javaClass:
			filePath = fs.GetJavaModelPath(s.workingDirectory, model.Location, model.LogicalName)

		case javaClassDep:
			filePath = fs.GetGenModelPath(s.workingDirectory, model.Location)
		case javaWar:
			compilerService := compiler.NewService(s.compilerServiceAddress)
			warFilePath, err := compilerService.CompileModel(
				s.workingDirectory,
				model.Location,
				model.LogicalName,
				"war",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filePath = warFilePath

		case javaJar:
			compilerService := compiler.NewService(s.compilerServiceAddress)
			jarFilePath, err := compilerService.CompileModel(
				s.workingDirectory,
				model.Location,
				model.LogicalName,
				"jar",
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filePath = jarFilePath
		}

		// Delegate to builtin.
		// Can result in 200, 404, 403 or 500 based on file availability and permissions.
		http.ServeFile(w, r, filePath)
		return

	default:
		http.Error(w, fmt.Sprintf("Invalid %s: %s", paramArtifact, artifact), http.StatusBadRequest)
		return
	}
}