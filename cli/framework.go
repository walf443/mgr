package main

type Framework interface {
	DetectFramework() bool
	GetSchema() (string, error)
	GetCurrentSchema() (string, error)
}

var SupportFrameworks []Framework

func init() {
	SupportFrameworks = append(SupportFrameworks, &RubyOnRailsFramework{})
}

func DetectFramework() Framework {
	for _, framework := range SupportFrameworks {
		if ( framework.DetectFramework() ) {
			return framework
		}
	}
	return nil
}
