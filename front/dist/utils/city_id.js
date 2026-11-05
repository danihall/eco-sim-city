const CITY_ID = "foss-eco-sim-city-hex";
function setCityID(id) {
    try {
        window.localStorage.setItem(CITY_ID, id.toString());
    }
    catch (reason) {
        console.error(reason);
    }
    return id;
}
function getCityID() {
    const id = window.localStorage.getItem(CITY_ID);
    return id?.length ? Number(id) : null;
}
export { setCityID, getCityID };
