# This is a simple web server to run just to test
# the basic functionality of the proxy. Not sure
# I'll keep it here forever, but we'll see.

import tornado.ioloop
import tornado.web


class MainHandler(tornado.web.RequestHandler):
    def get(self):
        print("MAIN HEADERS:\n\n", self.request.headers)
        self.write("MAIN!")


class FooHandler(tornado.web.RequestHandler):
    def get(self):
        print("FOO HEADERS:\n\n", self.request.headers)
        self.write("FOO!")


class BarHandler(tornado.web.RequestHandler):
    def get(self):
        print("BAR HEADERS:\n\n", self.request.headers)
        self.write("BAR!")


def make_app():
    return tornado.web.Application([
        (r"/", MainHandler),
        (r"/foo", FooHandler),
        (r"/bar.*", BarHandler),
    ])


if __name__ == "__main__":
    app = make_app()
    app.listen(8080)
    tornado.ioloop.IOLoop.current().start()

