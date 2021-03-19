package camunda

// FormVariableFilter is a filter for requesting form variables of proc
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-form-variables/#query-parameters
type FormVariableFilter struct {
	// VariableNames A comma-separated list of variable names.
	// Allows restricting the list of requested variables to the variable
	// names in the list.It is best practice to restrict the list of variables
	// to the variables actually required by the form in order to minimize
	// fetching of data.If the query parameter is ommitted all variables are
	// fetched.If the query parameter contains non-existent variable names,
	// the variable names are ignored.
	VariableNames string `query:"variableNames,omitempty"`

	// DeserializeValues Determines whether serializable variable values
	// (typically variables that store custom Java objects) should be deserialized
	// on server side (default true).
	//
	// If set to true, a serializable variable will be deserialized on server
	// side and transformed to JSON using Jackson's POJO/bean property introspection
	// feature. Note that this requires the Java classes of the variable value to
	// be on the REST API's classpath.
	//
	// If set to false, a serializable variable will be returned in its serialized
	// format. For example, a variable that is serialized as XML will be returned
	// as a JSON string containing XML.
	//
	// Note: While true is the default value for reasons of backward compatibility,
	// we recommend setting this parameter to false when developing web applications
	// that are independent of the Java process applications deployed to the engine.
	DeserializeValues bool `query:"deserializeValues,omitempty"`
}
