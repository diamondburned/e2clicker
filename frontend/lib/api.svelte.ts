export let userID = $state(localStorage.getItem("userID"));
export const isLoggedIn = () => !!userID;
