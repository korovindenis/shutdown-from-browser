// Request to backend
export default async function HandleClick(mode, TimeShutDown) {
  try {
    const response = await fetch('/api/v1/server-power/', {
      method: 'POST',
      body: JSON.stringify({ mode, TimeShutDown }),
      headers: {
        Accept: 'application/json',
      },
    });
      if (response.status === 200) {
        const result = await response.json();
        console.log(`Result is ${JSON.stringify(result, null, 4)}`);
      } else if (response.status !== 204) {
        throw new Error(`Error! status: ${response.status}`);
      }
  } catch (err) {
    console.error(err.message);
  }
};