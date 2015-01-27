package main

type Framework interface {
	DetectFramework() bool
	GetSchema() (string, error)
	GetCurrentSchema() (string, error)
}

var SupportingFrameworks []Framework

func init() {
	SupportingFrameworks = append(SupportingFrameworks, &RubyOnRailsFramework{})
}

func DetectFramework() Framework {
	for _, framework := range SupportingFrameworks {
		if ( framework.DetectFramework() ) {
			return framework
		}
	}
	return nil
}
