// Request to backend
export default async function HandleClick(mode, TimeShutDown=(new Date()).toISOString()){
    try {
      const response = await fetch('/api/v1/server-power/', {
        method: 'POST',
        body: JSON.stringify({ mode, TimeShutDown }),
        headers: {
          Accept: 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(`Error! status: ${response.status}`);
      }
  
      const result = await response.json();
      console.log('result is: ', JSON.stringify(result, null, 4));
    } catch (err) {
      console.error(err.message);
    }
  };