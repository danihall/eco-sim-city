import { getCityID } from "utils";

type DeleteCityPayload = {
  email: string;
  password: string;
  name: string;
  mayor_name: string;
};

async function deleteCity(payload: DeleteCityPayload): Promise<void> {
  const method = "DELETE";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}`;
  const body = JSON.stringify(payload);

  const response = await fetch(url, { method, headers, body });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }
}

export type { DeleteCityPayload };
export { deleteCity };
