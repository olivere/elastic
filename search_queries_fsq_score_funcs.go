// Copyright 2012-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"strings"
)

// ScoreFunction is used in combination with the Function Score Query.
type ScoreFunction interface {
	Name() string
	GetWeight() *float64
	Source() interface{}
}

// -- Exponential Decay --

type ExponentialDecayFunction struct {
	fieldName string
	origin    interface{}
	scale     interface{}
	decay     *float64
	offset    interface{}
	weight    *float64
}

func NewExponentialDecayFunction() ExponentialDecayFunction {
	return ExponentialDecayFunction{}
}

func (fn ExponentialDecayFunction) Name() string {
	return "exp"
}

func (fn ExponentialDecayFunction) FieldName(fieldName string) ExponentialDecayFunction {
	fn.fieldName = fieldName
	return fn
}

func (fn ExponentialDecayFunction) Origin(origin interface{}) ExponentialDecayFunction {
	fn.origin = origin
	return fn
}

func (fn ExponentialDecayFunction) Scale(scale interface{}) ExponentialDecayFunction {
	fn.scale = scale
	return fn
}

func (fn ExponentialDecayFunction) Decay(decay float64) ExponentialDecayFunction {
	fn.decay = &decay
	return fn
}

func (fn ExponentialDecayFunction) Offset(offset interface{}) ExponentialDecayFunction {
	fn.offset = offset
	return fn
}

func (fn ExponentialDecayFunction) Weight(weight float64) ExponentialDecayFunction {
	fn.weight = &weight
	return fn
}

func (fn ExponentialDecayFunction) GetWeight() *float64 {
	return fn.weight
}

func (fn ExponentialDecayFunction) Source() interface{} {
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source[fn.fieldName] = params
	if fn.origin != nil {
		params["origin"] = fn.origin
	}
	params["scale"] = fn.scale
	if fn.decay != nil && *fn.decay > 0 {
		params["decay"] = *fn.decay
	}
	if fn.offset != nil {
		params["offset"] = fn.offset
	}
	return source
}

// -- Gauss Decay --

type GaussDecayFunction struct {
	fieldName string
	origin    interface{}
	scale     interface{}
	decay     *float64
	offset    interface{}
	weight    *float64
}

func NewGaussDecayFunction() GaussDecayFunction {
	return GaussDecayFunction{}
}

func (fn GaussDecayFunction) Name() string {
	return "gauss"
}

func (fn GaussDecayFunction) FieldName(fieldName string) GaussDecayFunction {
	fn.fieldName = fieldName
	return fn
}

func (fn GaussDecayFunction) Origin(origin interface{}) GaussDecayFunction {
	fn.origin = origin
	return fn
}

func (fn GaussDecayFunction) Scale(scale interface{}) GaussDecayFunction {
	fn.scale = scale
	return fn
}

func (fn GaussDecayFunction) Decay(decay float64) GaussDecayFunction {
	fn.decay = &decay
	return fn
}

func (fn GaussDecayFunction) Offset(offset interface{}) GaussDecayFunction {
	fn.offset = offset
	return fn
}

func (fn GaussDecayFunction) Weight(weight float64) GaussDecayFunction {
	fn.weight = &weight
	return fn
}

func (fn GaussDecayFunction) GetWeight() *float64 {
	return fn.weight
}

func (fn GaussDecayFunction) Source() interface{} {
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source[fn.fieldName] = params
	if fn.origin != nil {
		params["origin"] = fn.origin
	}
	params["scale"] = fn.scale
	if fn.decay != nil && *fn.decay > 0 {
		params["decay"] = *fn.decay
	}
	if fn.offset != nil {
		params["offset"] = fn.offset
	}
	return source
}

// -- Linear Decay --

type LinearDecayFunction struct {
	fieldName string
	origin    interface{}
	scale     interface{}
	decay     *float64
	offset    interface{}
	weight    *float64
}

func NewLinearDecayFunction() LinearDecayFunction {
	return LinearDecayFunction{}
}

func (fn LinearDecayFunction) Name() string {
	return "linear"
}

func (fn LinearDecayFunction) FieldName(fieldName string) LinearDecayFunction {
	fn.fieldName = fieldName
	return fn
}

func (fn LinearDecayFunction) Origin(origin interface{}) LinearDecayFunction {
	fn.origin = origin
	return fn
}

func (fn LinearDecayFunction) Scale(scale interface{}) LinearDecayFunction {
	fn.scale = scale
	return fn
}

func (fn LinearDecayFunction) Decay(decay float64) LinearDecayFunction {
	fn.decay = &decay
	return fn
}

func (fn LinearDecayFunction) Offset(offset interface{}) LinearDecayFunction {
	fn.offset = offset
	return fn
}

func (fn LinearDecayFunction) Weight(weight float64) LinearDecayFunction {
	fn.weight = &weight
	return fn
}

func (fn LinearDecayFunction) GetWeight() *float64 {
	return fn.weight
}

func (fn LinearDecayFunction) Source() interface{} {
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source[fn.fieldName] = params
	if fn.origin != nil {
		params["origin"] = fn.origin
	}
	params["scale"] = fn.scale
	if fn.decay != nil && *fn.decay > 0 {
		params["decay"] = *fn.decay
	}
	if fn.offset != nil {
		params["offset"] = fn.offset
	}
	return source
}

// -- Script --

type ScriptFunction struct {
	script     string
	scriptFile string
	lang       string
	params     map[string]interface{}
	weight     *float64
}

func NewScriptFunction(script string) ScriptFunction {
	return ScriptFunction{
		script: script,
		params: make(map[string]interface{}),
	}
}

func (fn ScriptFunction) Name() string {
	return "script_score"
}

func (fn ScriptFunction) Script(script string) ScriptFunction {
	fn.script = script
	return fn
}

func (fn ScriptFunction) ScriptFile(scriptFile string) ScriptFunction {
	fn.scriptFile = scriptFile
	return fn
}

func (fn ScriptFunction) Lang(lang string) ScriptFunction {
	fn.lang = lang
	return fn
}

func (fn ScriptFunction) Param(name string, value interface{}) ScriptFunction {
	fn.params[name] = value
	return fn
}

func (fn ScriptFunction) Params(params map[string]interface{}) ScriptFunction {
	fn.params = params
	return fn
}

func (fn ScriptFunction) Weight(weight float64) ScriptFunction {
	fn.weight = &weight
	return fn
}

func (fn ScriptFunction) GetWeight() *float64 {
	return fn.weight
}

func (fn ScriptFunction) Source() interface{} {
	source := make(map[string]interface{})
	if fn.script != "" {
		source["script"] = fn.script
	}
	if fn.scriptFile != "" {
		source["script_file"] = fn.scriptFile
	}
	if fn.lang != "" {
		source["lang"] = fn.lang
	}
	if len(fn.params) > 0 {
		source["params"] = fn.params
	}
	return source
}

// -- Factor --

// FactorFunction is deprecated.
type FactorFunction struct {
	boostFactor *float32
}

func NewFactorFunction() FactorFunction {
	return FactorFunction{}
}

func (fn FactorFunction) Name() string {
	return "boost_factor"
}

func (fn FactorFunction) BoostFactor(boost float32) FactorFunction {
	fn.boostFactor = &boost
	return fn
}

func (fn FactorFunction) GetWeight() *float64 {
	return nil
}

func (fn FactorFunction) Source() interface{} {
	return fn.boostFactor
}

// -- Field value factor --

// FieldValueFactorFunction is a function score function that allows you
// to use a field from a document to influence the score.
// See http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-function-score-query.html#_field_value_factor.
type FieldValueFactorFunction struct {
	field    string
	factor   *float64
	missing  *float64
	weight   *float64
	modifier string
}

// NewFieldValueFactorFunction creates a new FieldValueFactorFunction.
func NewFieldValueFactorFunction() FieldValueFactorFunction {
	return FieldValueFactorFunction{}
}

// Name of the function score function.
func (fn FieldValueFactorFunction) Name() string {
	return "field_value_factor"
}

// Field is the field to be extracted from the document.
func (fn FieldValueFactorFunction) Field(field string) FieldValueFactorFunction {
	fn.field = field
	return fn
}

// Factor is the (optional) factor to multiply the field with. If you do not
// specify a factor, the default is 1.
func (fn FieldValueFactorFunction) Factor(factor float64) FieldValueFactorFunction {
	fn.factor = &factor
	return fn
}

// Modifier to apply to the field value. It can be one of: none, log, log1p,
// log2p, ln, ln1p, ln2p, square, sqrt, or reciprocal. Defaults to: none.
func (fn FieldValueFactorFunction) Modifier(modifier string) FieldValueFactorFunction {
	fn.modifier = modifier
	return fn
}

func (fn FieldValueFactorFunction) Weight(weight float64) FieldValueFactorFunction {
	fn.weight = &weight
	return fn
}

func (fn FieldValueFactorFunction) GetWeight() *float64 {
	return fn.weight
}

// Missing is used if a document does not have that field.
func (fn FieldValueFactorFunction) Missing(missing float64) FieldValueFactorFunction {
	fn.missing = &missing
	return fn
}

// Source returns the JSON to be serialized into the query.
func (fn FieldValueFactorFunction) Source() interface{} {
	source := make(map[string]interface{})
	if fn.field != "" {
		source["field"] = fn.field
	}
	if fn.factor != nil {
		source["factor"] = *fn.factor
	}
	if fn.missing != nil {
		source["missing"] = *fn.missing
	}
	if fn.modifier != "" {
		source["modifier"] = strings.ToLower(fn.modifier)
	}
	return source
}

// -- Weight Factor --

type WeightFactorFunction struct {
	weight float64
}

func NewWeightFactorFunction(weight float64) WeightFactorFunction {
	return WeightFactorFunction{weight: weight}
}

func (fn WeightFactorFunction) Name() string {
	return "weight"
}

func (fn WeightFactorFunction) Weight(weight float64) WeightFactorFunction {
	fn.weight = weight
	return fn
}

func (fn WeightFactorFunction) GetWeight() *float64 {
	return &fn.weight
}

func (fn WeightFactorFunction) Source() interface{} {
	return fn.weight
}

// -- Random --

type RandomFunction struct {
	seed   interface{}
	weight *float64
}

func NewRandomFunction() RandomFunction {
	return RandomFunction{}
}

func (fn RandomFunction) Name() string {
	return "random_score"
}

// Seed is documented in 1.6 as a numeric value. However, in the source code
// of the Java client, it also accepts strings. So we accept both here, too.
func (fn RandomFunction) Seed(seed interface{}) RandomFunction {
	fn.seed = seed
	return fn
}

func (fn RandomFunction) Weight(weight float64) RandomFunction {
	fn.weight = &weight
	return fn
}

func (fn RandomFunction) GetWeight() *float64 {
	return fn.weight
}

func (fn RandomFunction) Source() interface{} {
	source := make(map[string]interface{})
	if fn.seed != nil {
		source["seed"] = fn.seed
	}
	return source
}
