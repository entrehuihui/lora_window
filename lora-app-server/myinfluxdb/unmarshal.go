package myinfluxdb

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"

	"github.com/pkg/errors"
)

type alarmInfo struct {
	humidity_low     bool
	humidity_high    bool
	temperature_low  bool
	temperature_high bool
	electricity_low  bool
}

type humiture struct {
	temp    float64
	hum     float64
	elec    float64
	up_date time.Time
	alarm   byte
}

func unmarshalAlarm(alarm byte) alarmInfo {
	var info alarmInfo
	if alarm&0x01 == 0x01 {
		info.humidity_high = true
	}
	if alarm&0x02 == 0x02 {
		info.temperature_high = true
	}
	if alarm&0x04 == 0x04 {
		info.humidity_low = true
	}
	if alarm&0x08 == 0x08 {
		info.temperature_low = true
	}
	if alarm&0x10 == 0x10 {
		info.electricity_low = true
	}
	return info
}

func (a alarmInfo) String() string {
	buf := bytes.NewBufferString("")
	if a.humidity_high {
		buf.WriteString("湿度过高;")
	}
	if a.humidity_low {
		buf.WriteString("湿度过低;")
	}
	if a.temperature_high {
		buf.WriteString("温度过高;")
	}
	if a.temperature_low {
		buf.WriteString("温度过低;")
	}
	if a.electricity_low {
		buf.WriteString("电量过低;")
	}
	return string(buf.Bytes())
}

// 解码温湿度
//  frame_header function_code  timestamp  humidity     temperature  electricity  alarm crc  frame_tail
//      1             1             4         1          4(1+1+1+1)      1         1     1       1
func decodeHumiture(data []byte) ([]humiture, error) {
	if len(data) < 14 {
		return nil, errors.New("humiture data format error.")
	}
	count := int(math.Ceil(float64(len(data)) / 16.))

	if len(data)-(count-1)*16 < 14 {
		count = count - 1
	}
	results := make([]humiture, 0, count)
	for i := 0; i < count; i++ {
		var (
			start  int
			end    int
			result = humiture{}
		)
		if data[1+i*16] != 0x01 {
			continue
		}
		start = 2 + (i * 16)
		end = 6 + (i * 16)
		_up_date := binary.BigEndian.Uint32(data[start:end])
		result.up_date = time.Unix(int64(_up_date), 0)

		start = end
		end = end + 1
		result.hum = float64(data[start])

		start = end
		end = end + 1
		symbol := data[start]
		var flag float64 = 1.0
		if symbol == '-' {
			flag = -1.0
		}
		start = end
		end = end + 2
		temp_int := data[start]
		start = end
		end = end + 1
		_temp_dec := data[start]
		var temp_dec float64
		switch {
		case _temp_dec > 100:
			temp_dec = float64(_temp_dec) / 1000.
		case _temp_dec > 10:
			temp_dec = float64(_temp_dec) / 100.
		default:
			temp_dec = float64(_temp_dec) / 10.
		}
		result.temp = (float64(temp_int) + temp_dec) * flag
		start = end
		end = end + 1
		result.elec = float64(data[start])
		start = end
		end = end + 1
		result.alarm = data[start]
		results = append(results, result)
	}

	return results, nil
}
