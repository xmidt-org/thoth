package thoth

const mediaTypeKey = "_mediaType"

// Model is a template model.
type Model map[string]interface{}

type kvEntry struct {
	key   string
	value interface{}
}

// KeyValues is an immutable, efficient list of key/value pairs for
// enriching Models prior to template execution.
type KeyValues struct {
	entries []kvEntry
}

// Append adds key/value pairs to this list.  If more is not empty,
// a new, distinct KeyValues is created and returned with the Model's
// key/value pairs added.  If more is empty, this KeyValues is returned.
//
// Duplicate keys are allowed.  The final value of that key is simply
// the value last added to this list.
//
// This method does not modify this instance.
func (kvs KeyValues) Append(more Model) KeyValues {
	if len(more) == 0 {
		return kvs
	}

	merged := KeyValues{
		entries: make([]kvEntry, 0, len(kvs.entries)),
	}

	merged.entries = append(merged.entries, kvs.entries...)
	for k, v := range more {
		merged.entries = append(merged.entries, kvEntry{
			key:   k,
			value: v,
		})
	}

	return merged
}

// Extend adds the key/value pairs in another list to this list.  If more
// is empty, this instance is returned.  If this instance is empty, then
// more is returned.  Otherwise, a new, distinct KeyValues is returned that
// is the merge of the two.
//
// Duplicate keys are allowed.  The final value of that key is simply
// the value last added to this list.
//
// This method does not modify this instance.
func (kvs KeyValues) Extend(more KeyValues) KeyValues {
	if len(kvs.entries) == 0 {
		return more
	} else if len(more.entries) == 0 {
		return kvs
	}

	merged := KeyValues{
		entries: make([]kvEntry, 0, len(kvs.entries)+len(more.entries)),
	}

	merged.entries = append(merged.entries, kvs.entries...)
	merged.entries = append(merged.entries, more.entries...)
	return merged
}

// ApplyDefaults uses the key/value pairs in this list as defaults
// for the given model.  Any key present in this list but not present
// in the model will be set to the value in this list.  Existing keys
// in the model that are not in this list are left untouched.
func (kvs KeyValues) ApplyDefaults(m Model) {
	for _, e := range kvs.entries {
		if _, present := m[e.key]; !present {
			m[e.key] = e.value
		}
	}
}

// ApplyOverrides unconditionally sets each key in this list on the given model.
// Existing keys in the model that are in this list are overridden.
func (kvs KeyValues) ApplyOverrides(m Model) {
	for _, e := range kvs.entries {
		m[e.key] = e.value
	}
}
