package commands

import "fmt"

type Version struct {
	Label   string `json:"label"`
	Version string `json:"version"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Tag     string `json:"tag"`
	Risk    string `json:"risk"`
}

// Complete version list matching CC Version Guard database
// All URLs sourced from official ByteDance CDN
var versions = []Version{
	// 5.4.0 Beta series
	{Label: "5.4.0 (Beta6)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1991_beta6_capcutpc_beta_creatortool.exe", Type: "latest", Tag: "Latest", Risk: "High"},
	{Label: "5.4.0 (Beta5)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1988_beta5_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.4.0 (Beta4)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1982_beta4_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.4.0 (Beta3)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1979_beta3_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.4.0 (Beta2)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1978_beta2_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.4.0 (Beta1)", Version: "5.4.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1976_beta1_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	// 5.3.0 Stable & Beta series
	{Label: "5.3.0 (Latest)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1964_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.3.0 (Test2)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1961_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.3.0 (Test1)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1957_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.3.0 (Beta5)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1962_beta5_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta4)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1956_beta4_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta3)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1952_beta3_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta2 Latest)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1949_beta2_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta2 Test1)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1947_beta2_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta1 Latest)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1942_beta1_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	{Label: "5.3.0 (Beta1 Test1)", Version: "5.3.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1941_beta1_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	// 5.2.0 Stable & Beta series
	{Label: "5.2.0 (Latest)", Version: "5.2.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1950_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.2.0 (Test3)", Version: "5.2.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1946_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.2.0 (Test2)", Version: "5.2.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1940_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.2.0 (Test1)", Version: "5.2.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1939_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "High"},
	{Label: "5.2.0 (Beta8)", Version: "5.2.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1945_beta8_capcutpc_beta_creatortool.exe", Type: "beta", Tag: "Beta", Risk: "High"},
	// Older stable versions
	{Label: "4.8.0", Version: "4.8.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_4_8_0_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "Medium"},
	{Label: "4.5.0", Version: "4.5.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_4_5_0_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "Medium"},
	{Label: "3.5.0", Version: "3.5.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_3_5_0_capcutpc_0_creatortool.exe", Type: "stable", Tag: "Stable", Risk: "Low"},
	// First version
	{Label: "1.0.0 (Latest)", Version: "1.0.0", URL: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_1_0_0_44_capcutpc_0.exe", Type: "stable", Tag: "Stable", Risk: "Low"},
}

func GetVersions() []Version {
	return versions
}

func DownloadVersion(version string) error {
	var targetURL string
	for _, v := range versions {
		if v.Label == version {
			targetURL = v.URL
			break
		}
	}

	if targetURL == "" {
		return &DownloaderError{Message: fmt.Sprintf("version %q not found", version)}
	}

	return OpenURL(targetURL)
}

type DownloaderError struct {
	Message string
}

func (e *DownloaderError) Error() string {
	return e.Message
}
