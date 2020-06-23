package m_models

import (
	"fmt"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"strconv"
	"strings"
	"time"
)

//格式化的时候必须携带时区信息，使用 time.RFC3339 格式
//数据库里面所有的时间字段，必须使用该类型
type Time time.Time

func (this *Time) GobDecode(b []byte) error {
	t, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		return fmt.Errorf("m_models.Time.GobDecode: %w", err)
	}
	*this = Time(t)
	return nil
}

func (this Time) GobEncode() ([]byte, error) {
	return []byte(time.Time(this).Format(time.RFC3339)), nil
}

func (this Time) MarshalJSON() ([]byte, error) {
	t := time.Time(this)
	if t.IsZero() {
		return []byte(`"` + "" + `"`), nil
	}
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

func (this *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `""` || s == "" {
		s = "0001-01-01T08:00:00+08:00"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("m_models.Time.UnmarshalJSON: %w", err)
	}
	*this = Time(t)
	return nil
}

func (this Time) String() string {
	t := time.Time(this)
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func (this *Time) FromDB(b []byte) error {
	s := string(b)
	//如果字段是时间类型，sqlite3会按成时间对象，然后转成字节传递到当前方法
	if strings.Index(s, "T") == -1 {
		//sqlite3，没有做任何处理，这里我们按UTC时间解析，然后转成本地时间
		t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.UTC)
		if err != nil {
			return fmt.Errorf("m_models.Time.FromDB: %w", err)
		}
		*this = Time(t.Local())
	} else {
		//因为sqlite3的_loc参数配置为auto，所以sqlite3的默认时区是Local，所以sqlite3会将数据库的时间转成Local时区
		//此处已经是本地时区的字符串格式的时间，只需转成本地时区的时间对象即可
		t, err := time.ParseInLocation(time.RFC3339, s, time.Local)
		if err != nil {
			return fmt.Errorf("m_models.Time.FromDB T: %w", err)
		}
		*this = Time(t)
	}
	return nil
}

func (this *Time) ToDB() ([]byte, error) {
	//统一转成UTC时间入库
	t := time.Time(*this)
	if t.IsZero() {
		h_boot.Logger.Warning("m_models.Time.ToDB: zero time")
	}
	return []byte(t.UTC().Format("2006-01-02 15:04:05")), nil
}

func (this Time) Raw() time.Time {
	return time.Time(this)
}

//枚举类型
type Enum int

const EnumUNKNOWN = "UNKNOWN"

func (this *Enum) FromDB(b []byte) error {
	s := string(b)
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*this = Enum(i)
	return nil
}

func (this *Enum) ToDB() ([]byte, error) {
	return []byte(strconv.Itoa(int(*this))), nil
}
