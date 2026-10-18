# Bandwidth Hero Data Compression Service (Go Edition)

Welcome to the **Standalone Go** port of the Bandwidth Hero Data Compression Service 🚀. This service is designed to compress images on the fly, saving you bandwidth and improving your browsing experience.

Originally a Serverless Node.js application, this service has been fully rewritten in Go to be deployed as a highly efficient, single standalone binary. It compresses (optionally grayscale) given images to low-res [WebP](https://developers.google.com/speed/webp/) or JPEG format.

It downloads the original image and transforms it with `libvips` (via the `bimg` library) on the fly without saving images on disk.

**Benefits:**
- **Single Binary**: Easy to deploy without managing complex Node.js dependencies in production.
- **High Performance**: Leverages Go's concurrency and `libvips`'s blazing-fast C-based image processing.
- **Lower Resource Usage**: Optimized CPU and memory footprint compared to the previous Node.js serverless functions.

> Note: It downloads images on user's behalf (By passing in same headers to the domain with required image), passing cookies and user's IP address through to the origin host.

## Cloudflare Compatibility

This proxy forwards browser headers (User-Agent, Accept, Accept-Language, Accept-Encoding, Cookie, DNT, Referer) to make requests appear like legitimate browser requests. This helps avoid being blocked by Cloudflare CAPTCHA challenges that commonly trigger on bot-like requests.

## Deployment

Since this is a standalone Go web server, you will need a server (like a VPS) to run it.

### Prerequisites

You must install `libvips` on your system before compiling or running the proxy.
- **Ubuntu/Debian**: `sudo apt-get install libvips-dev pkg-config`
- **Alpine**: `apk add vips-dev pkgconfig`

### Building and Running

1. Clone the repository.
2. Build the binary using the provided Makefile:
   ```bash
   make build
   ```
3. Run the binary. By default, it listens on port `8080`, but you can configure this using the `PORT` environment variable:
   ```bash
   PORT=8080 ./bandwidth-hero-proxy
   ```

It is highly recommended to run this service behind a reverse proxy like Nginx, Caddy, or Cloudflare Tunnels.

Once deployed and running, open the **Data Compression Service** in your Bandwidth Hero extension, add your server URL (e.g., `http://your-server-ip:8080/api/index`), and you are good to go!
