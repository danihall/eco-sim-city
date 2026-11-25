type CreateCityPayload = {
  email: string;
  password: string;
  name: string;
  mayor_name: string;
};

async function createCity(payload: CreateCityPayload): Promise<number> {
  const method = "POST";
  const headers = { "Content-Type": "application/json" };
  const url = `/city`;
  const body = JSON.stringify(payload);

  const response = await fetch(url, { method, headers, body });

  if (!response.ok) {
    const reason = await response.text();
    throw new Error(`${response.status}\n${reason}`);
  }

  return response.json();
}

export type { CreateCityPayload };
export { createCity };
