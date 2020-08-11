package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSDL1(t *testing.T) {
	_, err := ParseSDL(`
		type TestExample1
			@stream(name: "test-example-1")
			@key(fields: ["a_aaa", "a_bbb"])
			@index(fields: ["a_bbb"], normalize: false)
		{
			a_aaa: String!
			a_bbb: Timestamp!
			a_ccc: [TestB!]
		}
		type TestB {
			b_aaa: Int
			b_bbb: Bytes
			b_ccc: TestC!
		}
		enum TestC {
			Aaa
			Bbb
			Ccc
		}
		type TestD {
			a_aaa: String!
			a_bbb: Timestamp!
			a_ccc: [TestExample1!]
		}
	`)
	assert.Nil(t, err)
}

func TestSDL3(t *testing.T) {
	_, err := ParseSDL(`
		type Test {
			a: Int!
		}
		type Test {
			b: Bool!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "has already been declared", err.Error())
}

func TestSDL4(t *testing.T) {
	_, err := ParseSDL(`
		type Test {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "no streams declared in input", err.Error())
}

func TestSDL5(t *testing.T) {
	_, err := ParseSDL(`
		type Test @stream {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "missing annotation '@key' with 'fields' arg in stream declaration at", err.Error())
}

func TestSDL6(t *testing.T) {
	_, err := ParseSDL(`
		type Test @stream(testA: "test", testB: "test") {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "unknown arg 'testA' for annotation '@stream'", err.Error())
}

func TestSDL9(t *testing.T) {
	_, err := ParseSDL(`
		type Test @stream @key(fields: [0, 1]) {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "arg 'fields' at .* is not a string or array of strings", err.Error())
}

func TestSDL10(t *testing.T) {
	_, err := ParseSDL(`
		type Test @stream @key(fields: ["a", "b"], external: whatever) {
			a: Int!
			b: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "parse error: .* unexpected \"whatever\" \\(expected .*\\)", err.Error())
}

func TestSDL11(t *testing.T) {
	_, err := ParseSDL(`
		type Test @stream @key(fields: true) {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "arg 'fields' at .* is not a string or array of strings", err.Error())
}

func TestSDL25(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @stream @key(fields: "a") {
			a: Int!
		}
		enum Bytes20 {
			Aa
			Bb
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "declaration of 'Bytes20' at .* overlaps with primitive type name", err.Error())
}

func TestSDL29(t *testing.T) {
	c, err := ParseSDL(`
		type TestA @stream(name: "test") @key(fields: "a") {
			a: Int!
			b: Int!
		}
	`)
	assert.Nil(t, err)
	assert.Equal(t, "test", c.Name)
}

func TestSDL30(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @stream(name: "test") @key(fields: "a") {
			a: Int!
		}
		type TestB @stream(name: "test") @key(fields: "a") {
			a: Int!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "found multiple objects in schema with '@stream' annotation", err.Error())
}

func TestSDL35(t *testing.T) {
	_, err := ParseSDL(`
		type TestA
			@stream
			@key(fields: "a")
			@hello
		{
			a: Int!
			b: String!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "unknown annotation '@hello' at ", err.Error())
}

func TestSDL36(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @key(fields: "a") {
			a: Int!
			b: String!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "cannot have '@key' or '@index' annotations on non-stream declaration at", err.Error())
}

func TestSDL37(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @stream(name: "") @key(fields: "a") {
			a: Int!
			b: String!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "stream arg 'name' at .* must be a non-empty string", err.Error())
}

func TestSDL39(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @stream @key(fields: "a", normalize: "true") {
			a: Int!
			b: String!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "key arg 'normalize' at .* must be a boolean", err.Error())
}

func TestSDL40(t *testing.T) {
	_, err := ParseSDL(`
		type TestA @stream @key(fields: "a", xxx: "true") {
			a: Int!
			b: String!
		}
	`)
	assert.NotNil(t, err)
	assert.Regexp(t, "unknown arg 'xxx' for annotation '@key' at", err.Error())
}
