package G

import (
	"bytes"
	"fmt"
)

type ServerMonitorErr struct {
	Msg string
}
func (this *ServerMonitorErr) Error() string {
	return this.Msg
}

func SerializeSetting(env string) string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("\"all\",%d,%d,%d", DefaultTpcTop[env], DefaultAlarmTop[env], DefaultAlarmInterval[env]))
	for name, value := range TpcTop[env] {
		buffer.WriteString(fmt.Sprintf(",\"%s\",%d,%d,%d", name, value, AlarmTop[env][name], AlarmInterval[env][name]))
	}
	return buffer.String()
}