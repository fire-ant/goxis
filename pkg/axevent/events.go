package axevent

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Cacsjep/goxis/pkg/axmanifest"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// KeyValueEntrie is a key-value pair for an AXEventKeyValueSet.
type KeyValueEntrie struct {
	Key       string
	Namespace *string
	Value     interface{}
	ValueType AXEventValueType
}

// Mark a key in the AXEventKeyValueSet as a source. A source key is an identifier used to distinguish between multiple instances of the same event declaration.
// E.g. if a device has multiple I/O ports then event declarations that represent the state of each port will have the same keys but different values.
// The key that represents which port the event represents should be marked as source and the key which represents the state should be marked as data.
// Please note that although it is possible to mark more than one key as a source, only events with zero or one source keys can be used to trigger actions.
type AxEventKeyValueSetSourceMark struct {
	Key       string
	Namespace *string
}

// Mark a key in the AXEventKeyValueSet as data. A data key is a key that represents the state of what the event represents.
// E.g. an event declaration that represents an I/O port should have a key marked as data which represents the state, high or low, of the port.
// Please note that although it is possible to mark more than one key as data, only events with one and only one data key can be used to trigger actions.
type AxEventKeyValueSetDataMark struct {
	Key       string
	Namespace *string
}

// Mark a key in AXEventKeyValueSet with an user defined tag.
type AxEventKeyValueSetUserDefineMark struct {
	Key       string
	Namespace *string
	Tag       *string
}

// Set the nice names of a key/value pair in the AXEventKeyValueSet.
// Nice names can be used to display human-readable information about the key/value pair.
type AxEventKeyValueSetNiceNames struct {
	Key           string
	Namespace     *string
	KeyNiceName   *string
	ValueNiceName *string
}

// NewTns1AxisEvent creates a new AXEventKeyValueSet with the given topics and key-value pairs.
func NewTns1AxisEvent(topic0 string, topic1 string, topic2 *string, topic3 *string, keyvalues []*KeyValueEntrie) (*AXEventKeyValueSet, error) {
	kvs := NewAXEventKeyValueSet()
	if err := kvs.AddKeyValue("topic0", &OnfivNameSpaceTns1, topic0, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic0: %w", err)
	}
	if err := kvs.AddKeyValue("topic1", &OnfivNameSpaceTnsAxis, topic1, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic1: %w", err)
	}
	if topic2 != nil {
		if err := kvs.AddKeyValue("topic2", &OnfivNameSpaceTnsAxis, *topic2, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if topic3 != nil {
		if err := kvs.AddKeyValue("topic3", &OnfivNameSpaceTnsAxis, *topic3, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if keyvalues != nil {
		for _, kv := range keyvalues {
			if err := kvs.AddKeyValue(kv.Key, kv.Namespace, kv.Value, kv.ValueType); err != nil {
				return nil, fmt.Errorf("failed to add key-value for %s: %w", kv.Key, err)
			}
		}

	}
	return kvs, nil
}

// TnsAxisEvent creates a new AXEventKeyValueSet with the given topics and key-value pairs.
func NewTnsAxisEvent(topic0 string, topic1 string, topic2 *string, topic3 *string, keyvalues []*KeyValueEntrie) (*AXEventKeyValueSet, error) {
	kvs := NewAXEventKeyValueSet()
	if err := kvs.AddKeyValue("topic0", &OnfivNameSpaceTnsAxis, topic0, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic0: %w", err)
	}
	if err := kvs.AddKeyValue("topic1", &OnfivNameSpaceTnsAxis, topic1, AXValueTypeString); err != nil {
		return nil, fmt.Errorf("failed to add key-value for topic1: %w", err)
	}
	if topic2 != nil {
		if err := kvs.AddKeyValue("topic2", &OnfivNameSpaceTnsAxis, *topic2, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
	}
	if topic3 != nil {
		if err := kvs.AddKeyValue("topic3", &OnfivNameSpaceTnsAxis, *topic3, AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic3: %w", err)
		}
	}
	if keyvalues != nil {
		for _, kv := range keyvalues {
			if err := kvs.AddKeyValue(kv.Key, kv.Namespace, kv.Value, kv.ValueType); err != nil {
				return nil, fmt.Errorf("failed to add key-value for %s: %w", kv.Key, err)
			}
		}

	}
	return kvs, nil
}

// UnmarshalEvent unmarshals the given event into the provided struct.
func UnmarshalEvent(e *Event, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		field := val.Elem().Field(i)
		if !field.CanSet() {
			continue
		}
		fieldType := val.Elem().Type().Field(i)

		key := fieldType.Tag.Get("eventKey")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}

		switch field.Kind() {
		case reflect.Int:
			if intValue, err := e.Kvs.GetInteger(key, nil); err == nil {
				field.SetInt(int64(intValue))
			} else {
				return fmt.Errorf("error getting integer for key %s: %v", key, err)
			}
		case reflect.Float64:
			if fValue, err := e.Kvs.GetDouble(key, nil); err == nil {
				field.SetFloat(fValue)
			} else {
				return fmt.Errorf("error getting double for key %s: %v", key, err)
			}
		case reflect.String:
			if sValue, err := e.Kvs.GetString(key, nil); err == nil {
				field.SetString(sValue)
			} else {
				return fmt.Errorf("error getting string for key %s: %v", key, err)
			}
		case reflect.Bool:
			if boolValue, err := e.Kvs.GetBoolean(key, nil); err == nil {
				field.SetBool(boolValue)
			} else {
				return fmt.Errorf("error getting boolean for key %s: %v", key, err)
			}
		}
	}

	return nil
}

func NewCameraApplicationPlatformEvent(
	app_setup axmanifest.Setup,
	event_name string,
	event_nice_name *string,
	kvs_entries []*KeyValueEntrie,
	source_markers []*AxEventKeyValueSetSourceMark,
	data_markers []*AxEventKeyValueSetDataMark,
	user_defined_markers []*AxEventKeyValueSetUserDefineMark,
	nice_names []*AxEventKeyValueSetNiceNames) (*AXEventKeyValueSet, error) {
	kvs, err := NewTnsAxisEvent(
		"CameraApplicationPlatform",
		app_setup.AppName,
		utils.NewStringPointer(event_name),
		nil,
		kvs_entries,
	)

	if err != nil {
		return nil, err
	}

	for _, source_marker := range source_markers {
		if err := kvs.MarkAsSource(source_marker.Key, source_marker.Namespace); err != nil {
			return nil, err
		}
	}

	for _, data_marker := range data_markers {
		if err := kvs.MarkAsData(data_marker.Key, data_marker.Namespace); err != nil {
			return nil, err
		}
	}

	for _, user_defined_marker := range user_defined_markers {
		if err := kvs.MarkAsUserDefined(user_defined_marker.Key, user_defined_marker.Namespace, user_defined_marker.Tag); err != nil {
			return nil, err
		}
	}

	var nice_name string
	if event_nice_name != nil {
		nice_name = fmt.Sprintf("%s: %s", app_setup.FriendlyName, *event_nice_name)
	} else {
		nice_name = fmt.Sprintf("%s: %s", app_setup.FriendlyName, event_name)
	}

	topic_nice_name := AxEventKeyValueSetNiceNames{
		Key: "topic2", Namespace: &OnfivNameSpaceTnsAxis, ValueNiceName: utils.NewStringPointer(nice_name),
	}
	nice_names = append(nice_names, &topic_nice_name)

	for _, nice_name := range nice_names {
		if err := kvs.AddNiceNames(nice_name.Key, nice_name.Namespace, nice_name.KeyNiceName, nice_name.ValueNiceName); err != nil {
			return nil, err
		}
	}

	return kvs, nil
}
