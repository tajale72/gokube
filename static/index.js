
 async function loadWeatherForecast() {
   try {
     const lat =
       document.getElementById("lat").value;


     const long =
       document.getElementById("long").value;


     const response = await fetch(
       `/weather?lat=${lat}&long=${long}`,
       {
         headers: {
           "X-API-KEY": "111"
         }
       }
     );


     if (!response.ok) {
       throw new Error(
         "Failed to fetch weather"
       );
     }


     const data = await response.json();


     // Updated response mapping
     document.getElementById(
       "shortForecast"
     ).textContent =
       data.short_forecast;


     document.getElementById(
       "temperature"
     ).textContent =
       `${data.temperature}°${data.temperature_unit}`;


     document.getElementById(
       "characterization"
     ).textContent =
       data.characterization;


   } catch (error) {
     console.error(error);


     document.getElementById(
       "shortForecast"
     ).textContent =
       "Weather unavailable";
   }
 }
