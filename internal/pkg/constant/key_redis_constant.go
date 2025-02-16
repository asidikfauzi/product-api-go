package constant

const (
	ByIDPattern = "-by-id-%s"

	PrefixCategory       = "category"
	AllCategoriesKey     = PrefixCategory + "-all-%s-%s"
	CategoryByIdKey      = PrefixCategory + ByIDPattern
	DeleteAllCategoryKey = PrefixCategory + "-all"

	PrefixMeasurement       = "measurement"
	AllMeasurementsKey      = PrefixMeasurement + "all-%s-%s"
	MeasurementByIdKey      = PrefixMeasurement + ByIDPattern
	DeleteAllMeasurementKey = PrefixMeasurement + "-all"

	PrefixProduct       = "product"
	AllProductsKey      = PrefixProduct + "-all-%s-%s-%s"
	ProductByIdKey      = PrefixProduct + ByIDPattern
	DeleteAllProductKey = PrefixProduct + "-all"
)
