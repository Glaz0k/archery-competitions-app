import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../model/series_model.dart';

class SeriesPage extends StatelessWidget {
  const SeriesPage({super.key});

  final String title = 'Серии';

  @override
  Widget build(BuildContext context) {
    return Consumer<SeriesModel>(
      builder:
          (context, seriesModel, _) => Scaffold(
            appBar: AppBar(
                backgroundColor: Theme.of(context).colorScheme.inversePrimary,
                title: Text(title)),
            floatingActionButton: FloatingActionButton(
              onPressed: seriesModel.addNewSeries,
              child: Icon(Icons.add),
            ),
            body: ListView(
              scrollDirection: Axis.vertical,
              children: [
                for (var (index, series) in seriesModel.seriesList.indexed)
                  SeriesCard(
                    name: Text("Серия ${index + 1}"),
                    series: series,
                    active: index + 1 == seriesModel.seriesList.length,
                    onAdd: () => seriesModel.addToLast(index),
                    onChange:
                        (index, score) =>
                            seriesModel.changeElementInLast(index, score),
                  ),
              ],
            ),
          ),
    );
  }
}

class SeriesCard extends StatefulWidget {
  final Widget name;
  final Series series;
  final bool active;
  final void Function() onAdd;
  final void Function(int, int) onChange;

  const SeriesCard({
    super.key,
    required this.name,
    required this.series,
    required this.onAdd,
    required this.onChange,
    required this.active,
  });

  @override
  State<StatefulWidget> createState() {
    return _SeriesCardState();
  }
}

class _SeriesCardState extends State<SeriesCard> {
  int? _selected;

  @override
  Widget build(BuildContext context) {
    List<Widget> body = [
      widget.name,
      ...widget.series.scores.indexed.map((r) {
        var (index, score) = r;
        var isSelected = index == _selected;
        var text = Text('$score');
        if (isSelected) {
          return FilledButton(
            onPressed: widget.active ? () => _selectScore(null) : null,
            child: text,
          );
        } else {
          return FilledButton.tonal(
            onPressed: widget.active ? () => _selectScore(index) : null,
            child: text,
          );
        }
      }),
      FloatingActionButton(onPressed: widget.onAdd, child: Icon(Icons.add)),
    ];
    var selected = _selected;

    return Card(
      child: Column(
        children: [
          //ListView(scrollDirection: Axis.horizontal, children: body),
          Row(children: body,),
          AnimatedContainer(
            duration: Duration(microseconds: 200),
            child:
                selected != null
                    ? Slider(
                      min: 0.0,
                      max: 10.0,
                      value: widget.series.scores[selected].toDouble(),
                      divisions: 10,
                      onChanged:
                          (score) => widget.onChange(selected, score.round()),
                    )
                    : SizedBox.shrink(),
          ),
        ],
      ),
    );
  }

  void _selectScore(int? index) => setState(() {
    _selected = index;
  });
}
