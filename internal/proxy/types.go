package proxy

type DockerVolume struct {
	Name    		string
	Mountpoint    	string
	Status			[]string
}

type PluginDescription struct {
	Implements		[]string
}