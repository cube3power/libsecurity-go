package password_test

import (
	"fmt"

	"github.com/ibm-security-innovation/libsecurity-go/password"
	"github.com/ibm-security-innovation/libsecurity-go/salt"
)

var (
	minPasswordLength = 1
	maxPasswordLength = 255
)

// Example of how to use the password.
// 1. Create a new password.
// 2. Verify that the initial password is set correctly
// 3. Change the user's password
// 4. Verify that the old password is not valid anymore
// 5. Verify that the new password is valid
// 6. Verify that the old password can't be used any more
//     (at least not as long as it remains in the old passwords list)
func ExampleUserPwd() {
	id := "User-1"
	pwd := []byte("a1B2c3d^@")
	saltStr, _ := salt.GetRandomSalt(8)

	userPwd, _ := password.NewUserPwd(pwd, saltStr, true)
	tPwd, _ := salt.GenerateSaltedPassword(pwd, minPasswordLength, maxPasswordLength, saltStr, -1)
	newPwd := password.GetHashedPwd(tPwd)
	err := userPwd.IsPasswordMatch(newPwd)
	if err != nil {
		fmt.Println("Error", err)
	}
	userNewPwd := []byte(string(pwd) + "a")
	newPwd, err = userPwd.UpdatePassword(userPwd.Password, userNewPwd, true)
	if err != nil {
		fmt.Printf("Password update for user %v to new password '%v' (%v) failed, error %v\n", id, newPwd, string(userNewPwd), err)
	} else {
		fmt.Printf("User '%v', updated password to '%v' (%v)\n", id, newPwd, string(userNewPwd))
	}
	err = userPwd.IsPasswordMatch(newPwd)
	if err != nil {
		fmt.Printf("Check of the new password, '%v' (%v), for user %v failed, error %v\n", newPwd, string(userNewPwd), id, err)
	} else {
		fmt.Printf("User '%v', new password '%v' (%v) verified successfully\n", id, newPwd, string(userNewPwd))
	}
	err = userPwd.IsPasswordMatch(pwd)
	if err == nil {
		fmt.Printf("Error: Old password '%v' (%v) for user %v accepted\n", pwd, string(pwd), id)
	} else {
		fmt.Printf("User '%v', Note that the old password '%v' (%v) cannot be used anymore\n", id, pwd, string(pwd))
	}
	newPwd, err = userPwd.UpdatePassword(userPwd.Password, pwd, true)
	if err == nil {
		fmt.Printf("Error: Password '%v' (typed password %v) for user %v was already used\n", newPwd, string(pwd), id)
	} else {
		fmt.Printf("Entity '%v'. Note that the old password (entered password) %v was already used\n", id, string(pwd))
	}
}

// Example of how to use the reset password function:
// This function resets the current password,
// selects a new password with short expiration time
// and lets the user use it exactly once
func ExampleUserPwd_ResetPassword() {
	id := "User1"
	pwd := []byte("a1b2C@3d4")

	saltStr, _ := salt.GetRandomSalt(10)
	userPwd, _ := password.NewUserPwd(pwd, saltStr, false)
	tmpPwd, _ := userPwd.ResetPassword()
	tPwd, _ := salt.GenerateSaltedPassword(tmpPwd, 1, 100, saltStr, -1)
	newPwd := password.GetHashedPwd(tPwd)
	err := userPwd.IsPasswordMatch(newPwd)
	if err != nil {
		fmt.Printf("Check of newly generated password '%v' for user %v failed, error %v\n", newPwd, id, err)
	} else {
		fmt.Printf("Entity %v, after resetting password '%v' verified successfully\n", id, newPwd)
	}
	err = userPwd.IsPasswordMatch(newPwd)
	if err == nil {
		fmt.Printf("Error: Newly generated password '%v' could be used only once\n", newPwd)
	} else {
		fmt.Printf("Newly generated password '%v', for entity %v, can only be used once\n", newPwd, id)
	}
}
