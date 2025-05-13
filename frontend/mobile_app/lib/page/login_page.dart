import 'package:flutter/material.dart';
import 'package:mobile_app/api/exceptions.dart';
import 'package:mobile_app/api/requests.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import '../api/api.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<StatefulWidget> createState() => _LoginPageState();
}


class _LoginPageState extends State<LoginPage> {
  Future<int> _userId = Future(() => 0);
  String? _errorMessage;
  final Credentials _credentials = Credentials("", "");
  @override
  Widget build(BuildContext context) {
    var api = context.watch<Api>();
    return Scaffold(
      appBar: OnionBar.withoutProfile("Вход", context),
      body: Form(
        child: Padding(
          padding: EdgeInsets.all(50),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              TextFormField(
                decoration: const InputDecoration(label: Text("Логин")),
                onSaved: (text) => _credentials.login = text ?? '',
                textInputAction: TextInputAction.next,
                forceErrorText: _errorMessage,
                onChanged: (_) => _clearError(),
                validator: _validate,
              ),
              TextFormField(
                decoration: const InputDecoration(label: Text("Пароль")),
                onSaved: (text) => _credentials.password = text ?? '',
                forceErrorText: _errorMessage,
                onChanged: (_) => _clearError(),
                validator: _validate,
              ),
              FutureBuilder(
                future: _userId,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.done) {
                    return FilledButton(
                      onPressed: () {
                        var formState = Form.of(context);
                        if (formState.validate()) {
                          formState.save();
                          setState(() {
                            _userId = api.login(_credentials).onError((e, st) {
                              var errorMessage = (e as OnionException).message;
                              setState(() {
                                _errorMessage = errorMessage;
                              });
                              return 0;
                            });
                          });
                        }
                      },
                      child: Text("Войти"),
                    );
                  } else {
                    return FilledButton(
                      onPressed: null,
                      child: Text("Входим..."),
                    );
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  String? _validate(String? value) {
    if (value == null || value.isEmpty) {
      return "Обязательное поле";
    }
    return null;
  }

  void _clearError() {
    if (_errorMessage != null) {
      setState(() {
        _errorMessage = null;
      });
    }
  }
}
