package proxy

type Volume struct {
	Name    		string
	Mountpoint    	string
	Status			[]string
}

type Capability struct {
	Scope			string
}

type VolumeDriverCapabilities struct {
	Capabilities	Capability
}

type PluginDescription struct {
	Implements		[]string
}