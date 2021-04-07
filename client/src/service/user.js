import axios from "axios";
import { base_url } from ".";

export async function registerUser(username, password, role) {
  var user = {
    username: username,
    password: password,
    role: role,
  };
  try {
    await axios.post(base_url + "register", user, {});
    return {
      success: true,
      error: null,
    };
  } catch (err) {
    return {
      success: false,
      error: err,
    };
  }
}

export async function loginUser(username, password) {
  var user = {
    username: username,
    password: password,
  };
  try {
    var response = await axios.post(base_url + "login", user, {});
    var token = response.data.Data;
    window.localStorage.setItem("token", token);
    return {
      success: true,
    };
  } catch (err) {
    return {
      success: false,
      error: err,
    };
  }
}
