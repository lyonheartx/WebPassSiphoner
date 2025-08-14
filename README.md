Chrome Password Capture Extension (Proof of Concept)

This Chrome extension demonstrates how passwords that are stored in password managers or typed in real time in the browser can be captured and sent to a server. It works with password managers such as LastPass, Keeper, Google Password Manager, and 1Password by reading autofilled passwords, avoiding the need to decrypt or breach cloud platforms. Captured credentials are sent to a web server, which can optionally forward them to a reverse shell listener. Passwords are also displayed in the Chrome console for testing purposes.

Features

Version 1: Simple web server that saved captured passwords to a .txt file.

Version 2: Full reverse shell capabilities, receiving passwords directly on the listener.

Note: A legitimate web server with certificates is not required for a local proof-of-concept web server, but it's possible to use a real web server to receive passwords. Chrome typically blocks untrusted certificates and https to http cross site POST requests. 

Why Golang?

Originally, a Python-based server was used. The project was converted to Go for:

Easier distribution and compilation into a standalone executable.

Firewall prompts are minimized; even if blocked, the server can still receive passwords locally.

Why a Web Server?

Chrome browsers are sandboxed, making it difficult to export data directly from the browser to the host system. A web server provides a straightforward way to receive captured data as it's the nature of javascript based extentions. Once the data is shared with the your web server, we control the web server and therefor the data. We can not output it to a file or siphon the credentials and send them through a reverse shell with version 2's reverse shell server. 

Usage

Start the server:

go run ReverseServer.go


or run the compiled executable:

./server


Start the listener (for reverse shell and live password output):

nc -lvnp 50000


Install the Chrome extension in Developer Mode or via enterprise deployment. There are methods using local group policies and the registry, as well as manipulating the encryption key to the local preferences of the chrome user, or more obvious; automated scripts that manipulate the browser to enable developer mode. It is difficult to get a rogue extention onto Chrome but not impossible. 

Monitor credentials: Captured passwords will be sent to the server and displayed in the listener terminal.

I kept the technical explanation intact but removed some wording that could be interpreted as promoting malicious activity. Since this is going on GitHub, you might also want to add a Disclaimer section stating that this is for educational or security research purposes only.

I can draft that disclaimer for you too—it would make the README safer and more professional. Do you want me to do that?
