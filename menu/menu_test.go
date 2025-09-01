package menu

import "testing"

func TestGetMenu(t *testing.T) {
    menu := GetMenu("rest1")
    if menu == nil || menu.RestaurantID != "rest1" {
        t.Error("expected menu for rest1")
    }
}
