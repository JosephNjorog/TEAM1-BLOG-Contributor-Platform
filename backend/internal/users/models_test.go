package users

import "testing"

func TestRole_Valid(t *testing.T) {
	valid := []Role{RoleSuperAdmin, RoleModerator, RoleGraphicDesigner, RolePublisher, RoleContributor}
	for _, r := range valid {
		if !r.Valid() {
			t.Errorf("%s.Valid() = false, want true", r)
		}
	}

	invalid := []Role{"", "admin", "Contributor", "super-admin"}
	for _, r := range invalid {
		if r.Valid() {
			t.Errorf("%s.Valid() = true, want false", r)
		}
	}
}
