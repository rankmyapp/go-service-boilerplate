package handlers

type UserRoutePermissions struct {
	Create []int
	List   []int
	Get    []int
	Update []int
	Delete []int
}

type ExportRoutePermissions struct {
	Create []int
}
