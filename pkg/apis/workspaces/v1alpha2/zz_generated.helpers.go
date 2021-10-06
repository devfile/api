package v1alpha2

//IsDefault returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func IsDefault(in *CommandGroup) bool {
	if in.IsDefault != nil {
		return *in.IsDefault
	}
	return false
}

//HotReloadCapable returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func HotReloadCapable(in *ExecCommand) bool {
	if in.HotReloadCapable != nil {
		return *in.HotReloadCapable
	}
	return false
}

//Parallel returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func Parallel(in *CompositeCommand) bool {
	if in.Parallel != nil {
		return *in.Parallel
	}
	return false
}

//MountSources returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func MountSources(in *Container) bool {
	if in.MountSources != nil {
		return *in.MountSources
	} else {
		if DedicatedPod(in) {
			return false
		}
		return true
	}
}

//DedicatedPod returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func DedicatedPod(in *Container) bool {
	if in.DedicatedPod != nil {
		return *in.DedicatedPod
	}
	return false
}

//RootRequired returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func RootRequired(in *Dockerfile) bool {
	if in.RootRequired != nil {
		return *in.RootRequired
	}
	return false
}

//Ephemeral returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func Ephemeral(in *Volume) bool {
	if in.Ephemeral != nil {
		return *in.Ephemeral
	}
	return false
}

//Secure returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func Secure(in *Endpoint) bool {
	if in.Secure != nil {
		return *in.Secure
	}
	return false
}
