package options

type MirrorOptions struct {
	Generic        *GenericOptions
	Etcd           *EtcdOptions
	Transport      *TransportOptions
	KeyValue       *KeyValueOptions
	RestfulServing *RestfulServingOptions
}

func NewMirrorOptions() *MirrorOptions {
	return &MirrorOptions{
		Generic:        NewGenericOptions(),
		Etcd:           NewEtcdOptions(),
		Transport:      NewTransportOptions(),
		KeyValue:       NewKeyValueOptions(),
		RestfulServing: NewDaemonOptions(),
	}
}

func (m *MirrorOptions) Flags() (sfs SectionFlagSet) {
	m.Generic.AddFlags(sfs.FlagSet("Generic flags"))
	m.Etcd.AddFlags(sfs.FlagSet("Etcd flags"))
	m.Transport.AddFlags(sfs.FlagSet("Transport flags"))
	m.KeyValue.AddFlags(sfs.FlagSet("KeyValue flags"))
	m.RestfulServing.AddFlags(sfs.FlagSet("Daemon flags"))
	return sfs
}

func (m *MirrorOptions) Validation() []error {
	var errs []error

	errs = append(errs, m.Generic.Validation()...)
	errs = append(errs, m.Etcd.Validate()...)
	errs = append(errs, m.Transport.Validate()...)
	errs = append(errs, m.KeyValue.Validate()...)
	errs = append(errs, m.RestfulServing.Validate()...)
	return errs
}
