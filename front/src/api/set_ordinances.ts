import { getCityID } from "utils";

type SetOrdinancePayload = {
  ordinance_names: OrdinanceType[];
};

async function setOrdinances(form_data: FormData): Promise<void> {
  const method = "PUT";
  const headers = { "Content-Type": "application/json" };
  const url = `/city/${getCityID()}/ordinances/set`;

  const payload: SetOrdinancePayload = {
    ordinance_names: form_data.keys().toArray() as OrdinanceType[],
  };
  const body = JSON.stringify(payload);

  const response = await fetch(url, { method, headers, body });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }
}

export { setOrdinances };
