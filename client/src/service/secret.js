import axios from "axios";
import { base_api_url } from ".";

export async function getSecrets() {
  let token = window.localStorage.getItem("token");
  try {
    var response = await axios.get(base_api_url + "secrets", {
      headers: {
        Authorization: "Bearer " + token,
      },
    });
    var data = response.data.Data.data;
    var secrets = [];
    if (data !== null) {
      data.forEach((d) => {
        secrets.push({
          alias: d.Alias,
          value: d.Value,
          role: d.Role,
          show: false,
          edit: false,
        });
      });
    }

    var deadNodes = [];
    if (response.data.Data.deadNodes !== null) {
      deadNodes = response.data.Data.deadNodes;
    }
    return {
      success: true,
      error: null,
      data: secrets,
      deadNodes: deadNodes,
    };
  } catch (err) {
    return {
      success: false,
      error: err,
    };
  }
}

export async function putSecret(alias, value, role) {
  let token = window.localStorage.getItem("token");
  try {
    await axios.put(
      base_api_url + "secret",
      {
        Alias: alias,
        Value: value,
        Role: Number(role),
      },
      {
        headers: {
          Authorization: "Bearer " + token,
        },
      }
    );
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

export async function deleteSecret(key) {
  let token = window.localStorage.getItem("token");
  try {
    await axios.delete(base_api_url + `secret?alias=${key}`, {
      headers: {
        Authorization: "Bearer " + token,
      },
    });
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
