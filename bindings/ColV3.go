/*
  Copyright (C) 2016 H2O.ai, Inc. <http://h2o.ai/>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as
  published by the Free Software Foundation, either version 3 of the
  License, or (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package bindings

type ColV3 struct {
	*Schema
	/** label */
	Label string `json:"label"`

	/** missing */
	MissingCount int64 `json:"missing_count"`

	/** zeros */
	ZeroCount int64 `json:"zero_count"`

	/** positive infinities */
	PositiveInfinityCount int64 `json:"positive_infinity_count"`

	/** negative infinities */
	NegativeInfinityCount int64 `json:"negative_infinity_count"`

	/** mins */
	Mins []float64 `json:"mins"`

	/** maxs */
	Maxs []float64 `json:"maxs"`

	/** mean */
	Mean float64 `json:"mean"`

	/** sigma */
	Sigma float64 `json:"sigma"`

	/** datatype: {enum, string, int, real, time, uuid} */
	Type string `json:"type"`

	/** domain; not-null for categorical columns only */
	Domain []string `json:"domain"`

	/** cardinality of this column's domain; not-null for categorical columns only */
	DomainCardinality int32 `json:"domain_cardinality"`

	/** data */
	Data []float64 `json:"data"`

	/** string data */
	StringData []string `json:"string_data"`

	/** decimal precision, -1 for all digits */
	Precision byte `json:"precision"`

	/** Histogram bins; null if not computed */
	HistogramBins []int64 `json:"histogram_bins"`

	/** Start of histogram bin zero */
	HistogramBase float64 `json:"histogram_base"`

	/** Stride per bin */
	HistogramStride float64 `json:"histogram_stride"`

	/** Percentile values, matching the default percentiles */
	Percentiles []float64 `json:"percentiles"`
}

func NewColV3() *ColV3 {
	return &ColV3{
		Label:                 "",
		MissingCount:          0,
		ZeroCount:             0,
		PositiveInfinityCount: 0,
		NegativeInfinityCount: 0,
		Mins:              nil,
		Maxs:              nil,
		Mean:              0.0,
		Sigma:             0.0,
		Type:              "",
		Domain:            nil,
		DomainCardinality: 0,
		Data:              nil,
		StringData:        nil,
		Precision:         0,
		HistogramBins:     nil,
		HistogramBase:     0.0,
		HistogramStride:   0.0,
		Percentiles:       nil,
	}
}
