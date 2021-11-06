package main

type Dictionary map[string]string
type DictionaryError string

const (
	ErrKeyNotFound = DictionaryError("word not found")
	ErrKeyExisted = DictionaryError("word existed")
)

// implement the error interface
func (d DictionaryError) Error() string {
	return string(d)
}

func (d Dictionary) Search(key string) (string, error) {
	value, ok := d[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (d Dictionary) Add(key string, value string) error {
	_, err := d.Search(key)
	switch err {
	case ErrKeyNotFound:
		d[key] = value
	case nil:
		return ErrKeyExisted
	default:
		return err
	}
	return nil
}

func (d Dictionary) Update(key string, newValue string) error {
	_, err := d.Search(key)
	switch err {
	case nil:
		d[key] = newValue
	case ErrKeyNotFound:
		return err
	default:
		return err
	}
	return nil
}

func (d Dictionary) Delete(key string) error {
	_, err := d.Search(key)
	if err == nil {
		delete(d, key)
		return nil
	}
	return err
}


