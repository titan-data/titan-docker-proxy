package proxy

type Volume struct {
	Name       string
	Mountpoint string
	Status     map[string]string
}

type Capability struct {
	Scope string
}

type VolumeCapabilities struct {
	Capabilities Capability
}

type VolumeRequest struct {
	Name string
}

type VolumeResponse struct {
	Err string
}

type ListVolumeResponse struct {
	Err     string
	Volumes []Volume
}

type GetVolumeResponse struct {
	Err    string
	Volume Volume
}

type GetPathResponse struct {
	Err        string
	Mountpoint string
}

type PluginDescription struct {
	Implements []string
}

type CreateVolumeRequest struct {
	Name string
	Opts map[string]interface{}
}
