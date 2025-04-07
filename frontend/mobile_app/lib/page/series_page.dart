import 'package:flutter/material.dart';
import 'package:mobile_app/model/series.dart';

class SeriesPage extends StatefulWidget {
  const SeriesPage({super.key});

  final String title = 'Серии';

  @override
  State<SeriesPage> createState() => _SeriesPageState();
}

class _SeriesPageState extends State<SeriesPage> {
  final List<Series> results = [Series()];
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Placeholder(),
      floatingActionButton: FloatingActionButton(
        onPressed: _addNewSeries,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ),
    );
  }
  void _addNewSeries() => setState(() {
    if (results.last.length > 5) {
      results.add(Series());
    } else {
     // int? newResult = _showNumberPicker();
    }
  });
}
