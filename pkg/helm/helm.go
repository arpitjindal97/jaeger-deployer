package helm

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"jaeger-tenant/pkg/structures"
	"jaeger-tenant/utils"
	"net/http"
	"os"
	"path"
)

func generateHelmTemplate(customerName, helmValues string) (string, error) {
	chart, _ := loader.Load("jaeger-0.19.8.tgz")

	kubeConfigPath := ""
	releaseName := customerName
	releaseNamespace := "default"
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(

		kube.GetConfig(kubeConfigPath, "", releaseNamespace),
		releaseNamespace,
		os.Getenv("HELM_DRIVER"),
		func(format string, v ...interface{}) {
			_ = fmt.Sprintf(format, v)
		}); err != nil {
		return "", err
	}

	iCli := action.NewInstall(actionConfig)
	iCli.Namespace = releaseNamespace
	iCli.ReleaseName = releaseName
	iCli.DryRun = true
	iCli.ClientOnly = true

	m := make(map[string]interface{})

	err := yaml.Unmarshal([]byte(helmValues), &m)
	if err != nil {
		return "", err
	}

	rel, err := iCli.Run(chart, m)
	if err != nil {
		return "", err
	}
	return rel.Manifest, nil
}

// GetTenantYaml generates the final tenant with helm and returns it
func GetTenantYaml(tenant *structures.Tenant, payload *structures.TenantPayload) (string, error) {
	values, err := utils.GetValuesYaml(tenant, payload)
	if err != nil {
		return "", err
	}
	helmValues, err := generateHelmTemplate(tenant.Customer, values)

	return helmValues, err
}

// DownloadHelmChart downloads the chart in current directory
func DownloadHelmChart() {

	option := getter.WithBasicAuth("", "")

	httpGetter, _ := getter.NewHTTPGetter(option)

	constructor := func(options ...getter.Option) (getter.Getter, error) {
		result := httpGetter
		return result, nil
	}

	provider := getter.Provider{
		Schemes: []string{"https"},
		New:     constructor,
	}

	url, _ := repo.FindChartInRepoURL("https://jaegertracing.github.io/helm-charts",
		"jaeger", "0.19.8", "", "", "", getter.Providers{provider})

	/*
		chartRepository,_ := repo.NewChartRepository(&repo.Entry{
			Name:     "",
			URL:      "https://jaegertracing.github.io/helm-charts",
			Username: "",
			Password: "",
			CertFile: "",
			KeyFile:  "",
			CAFile:   "",
		},getter.Providers{provider})
		indexPath,_ :=chartRepository.DownloadIndexFile()
		chartRepository.Load()

		indexFile,_ := repo.LoadIndexFile(indexPath)

		url:= indexFile.Entries["jaeger"][0].URLs[0]

	*/

	filename := path.Base(url)
	fmt.Println("Downloading ", url, " to ", filename)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	/*
		chartPathOptions := action.ChartPathOptions{
			RepoURL: "https://jaegertracing.github.io/helm-charts",
		}

		envSetting := cli.EnvSettings{
			KubeConfig:       "/Users/I353342/workspace/kube.yml",
			KubeContext:      "",
			Debug:            false,
			RegistryConfig:   "",
			RepositoryConfig: "",
			RepositoryCache:  "",
			PluginsDirectory: "",
		}

		str,err := chartPathOptions.LocateChart("jaeger",&envSetting)
		if err!= nil {
			panic(err)
		}
		fmt.Println(str)
	*/

}
