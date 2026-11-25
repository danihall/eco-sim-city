import { getCityID } from "utils";

async function getAllOrdinances(): Promise<Ordinance[]> {
  const method = "GET";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/ordinances`;

  const response = await fetch(url, { method, headers });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }

  return response.json();
}

export { getAllOrdinances };
