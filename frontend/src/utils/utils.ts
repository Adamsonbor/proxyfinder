export function SetCookie(name: string, value: string, expireSecond: number) {
	const expires = "expires=" + new Date(Date.now() + expireSecond * 1000).toUTCString();
	document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

export function GetCookie(name: string): string | null {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for (var i = 0; i < ca.length; i++) {
		var c = ca[i];
		if (c.trim().indexOf(nameEQ) == 0) {
			return c.substring(nameEQ.length, c.length);
		}
	}
	return null;
}

