import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/User.dart';
import 'package:mobile_app/page/widgets/appbar_widget.dart';
import 'package:mobile_app/page/widgets/profile_widget.dart';
import 'package:mobile_app/page/widgets/text_box.dart';

void main() {
  runApp(const ProfilePage());
}

class Onion extends StatelessWidget {
  const Onion({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});
  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text('You have pushed the button this many times:'),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ),
    );
  }
}

class AuthPage extends StatelessWidget {
  const AuthPage({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
        primarySwatch: Colors.green,
        scaffoldBackgroundColor: Colors.grey
      ),
      home: LoginPage()
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  _LoginPage createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {

  final _sizeBlackText = const TextStyle(fontSize: 12, color: Colors.black);
  final _sizeWhiteText = const TextStyle(fontSize: 15, color: Colors.white);

  final myController = TextEditingController();
  bool _submitted = false;
  var text = '';

  void _submit() {
    setState(() => _submitted = true);
    if (_errorText == null) {
      text = myController.text;
    }
  }

  @override
  void dispose() {
    // Clean up the controller when the widget is disposed.
    myController.dispose();
    super.dispose();
  }

  String? get _errorText {
    final text = myController.value.text;
    if (text.isEmpty) {
      return 'Поле не может быть пустым';
    }
    if (!text.contains('@')) {
      return 'Адрес должен содержать @';
    }
    return null;
  }

  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Вход'),
        backgroundColor: Colors.green,
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            SizedBox(
              width: 300,
              child: TextFormField(
                controller: myController,
                decoration: InputDecoration(
                  labelText: "Email",
                  errorText: _submitted ? _errorText : null,
                  icon: const Icon(Icons.person),
                  hintText: 'адрес электронной почты',

                ),
                validator: (String? value) {
                  return (value != null && !value.contains('@')) ? "Нужно использовать значок @" : null;
                },
                onChanged: null,
                keyboardType: TextInputType.emailAddress,
                style: _sizeBlackText,
              ),
            ),
            Container(
              width: 300,
              padding: EdgeInsets.only(top: 10.0),
              child: TextFormField(
                decoration: InputDecoration(
                    labelText: "Password",
                    hintText: 'пароль',
                    icon: Icon(Icons.person)
                ),
                obscureText: true,
                style: _sizeBlackText,
              ),
            ),
            Padding(
              padding:EdgeInsets.only(top: 25.0),
              child: ElevatedButton(
                  onPressed: null,
                  style: ButtonStyle(
                    minimumSize: WidgetStateProperty.all(const Size(200, 40)),
                    backgroundColor: WidgetStateProperty.all(Colors.green),
                    overlayColor: WidgetStateProperty.all(Colors.blueAccent)
                  ),
                  child: Text(
                    "Войти",
                    style: _sizeWhiteText,
                  )


              ),
            )
          ],
        ),
      ),
    );
  }
  
}

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  _ProfilePageState createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final user = UserPreferences.myUser;
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          primaryColor: Colors.green.shade300
        ),
        home: Scaffold(appBar: buildAppBar(context),
          body: ListView(
            physics: BouncingScrollPhysics(),
            children: [
              ProfileWidget(imagePath: user.imagePath, onClicked: () async {}),
              MyTextBox(text: user.name, sectionName: "Name"),
            ],)
        )
    );
  }
}

class EditProfilePage extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _EditProfilePage();
}

class _EditProfilePage extends State<EditProfilePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(

    );
  }
}