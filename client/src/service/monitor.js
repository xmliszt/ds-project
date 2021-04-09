import axios from "axios";
import { base_api_url } from ".";

// Getting monitoring information for display
export async function getMonitorInfo() {
  let token = window.localStorage.getItem("token");
  try {
    var response = await axios.get(base_api_url + "monitor", {
      headers: {
        Authorization: "Bearer " + token,
      },
    });
    return {
      success: true,
      data: response.data.Data,
    };
  } catch (err) {
    return {
      success: false,
      error: err,
    };
  }
}
