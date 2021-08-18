import 'main.dart';

Future<int> main() async {
  RWebClient c = RWebClient(serverAddress: "ws://127.0.0.1:1111/t");
  if (!await c.connect()) {
    print("error");
    return 0;
  }
  c.bind("test", (reply, arguments) {
    print(arguments);
    reply("Hi!");
  });
  print("Go");
  var ans = await c.call("test", List.empty());
  print("Called" + ".the response is " + ans);
  return 0;
}
