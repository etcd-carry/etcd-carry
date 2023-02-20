package options

type MirrorOptions struct {
	Generic   *GenericOptions
	Etcd      *EtcdOptions
	Transport *TransportOptions
}

func NewMirrorOptions() *MirrorOptions {
	return &MirrorOptions{
		Generic:   NewGenericOptions(),
		Etcd:      NewEtcdOptions(),
		Transport: NewTransportOptions(),
	}
}

func (m *MirrorOptions) Flags() (sfs SectionFlagSet) {
	m.Generic.AddFlags(sfs.FlagSet("Generic flags"))
	m.Etcd.AddFlags(sfs.FlagSet("Etcd flags"))
	m.Transport.AddFlags(sfs.FlagSet("Transport flags"))
	return sfs
}

func (m *MirrorOptions) Validation() []error {
	var errs []error

	errs = append(errs, m.Generic.Validation()...)
	errs = append(errs, m.Etcd.Validate()...)
	errs = append(errs, m.Transport.Validate()...)
	return errs
}
