import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import '../model/range_model.dart';

class SeriesPage extends StatelessWidget {
  const SeriesPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<RangeModel>(
      builder:
          (context, rangeModel, _) => Scaffold(
            appBar: OnionBar("Серии", context),
            floatingActionButton: FloatingActionButton.extended(
              onPressed: rangeModel.pushRange,
              icon: Icon(Icons.add),
              label: Text('Новая серия'),
            ),
            body: ListView(
              scrollDirection: Axis.vertical,
              children: [
                for (var (index, range) in rangeModel.ranges.indexed)
                  RangeCard(
                    name: Text("Серия ${index + 1}"),
                    range: range,
                    isActive: index == rangeModel.ranges.length - 1,
                  ),
              ],
            ),
          ),
    );
  }
}

class RangeCard extends StatefulWidget {
  final Widget name;
  final Range range;
  final bool isActive;

  const RangeCard({
    super.key,
    required this.name,
    required this.range,
    required this.isActive,
  });

  @override
  State<RangeCard> createState() => _RangeCardState();
}

class _RangeCardState extends State<RangeCard> {
  int? _selected;

  @override
  Widget build(BuildContext context) {
    var selected = _selected;
    var selectedShot = selected == null ? null : widget.range.shots[selected];
    var model = Provider.of<RangeModel>(context, listen: false);
    var maxShots = model.maxShots;
    var isShort = model.isShort;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                SizedBox(
                  width: 80,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      widget.name,
                      ListenableBuilder(
                        listenable: widget.range.score,
                        builder:
                            (context, _) =>
                                Text("Счёт: ${widget.range.score.value}/30"),
                      ),
                    ],
                  ),
                ),
                ListenableBuilder(
                  listenable: widget.range,
                  builder:
                      (context, _) => Expanded(
                        child: SizedBox(
                          height: 50,
                          child: ListView(
                            scrollDirection: Axis.horizontal,
                            children: [
                              for (var (index, shot)
                                  in widget.range.shots.indexed)
                                RangeCardButton(
                                  shot: shot,
                                  isSelected:
                                      index == _selected && widget.isActive,
                                  onPressed:
                                      widget.isActive
                                          ? () => setState(() {
                                            if (index == _selected) {
                                              _selected = null;
                                            } else {
                                              _selected = index;
                                            }
                                          })
                                          : null,
                                ),
                              if (widget.range.shots.length < maxShots &&
                                  widget.isActive)
                                IconButton(
                                  onPressed: () {
                                    widget.range.pushShot(Shot(0));
                                    setState(() {
                                      _selected = widget.range.shots.length - 1;
                                    });
                                  },
                                  icon: Icon(Icons.add),
                                ),
                            ],
                          ),
                        ),
                      ),
                ),
              ],
            ),
            if (selectedShot != null && widget.isActive)
              ListenableBuilder(
                listenable: selectedShot,
                builder: (context, _) {
                  return Slider(
                    value: selectedShot.value.rawValue.toDouble(),
                    min: isShort ? 5 : 0,
                    max: 11,
                    divisions: (isShort ? 11 - 5 : 11),
                    onChanged: (value) {
                      widget.range.changeShot(
                        selected!,
                        Shot.fromSlider(value, isShort),
                      );
                    },
                  );
                },
              ),
          ],
        ),
      ),
    );
  }
}

class RangeCardButton extends StatelessWidget {
  final ValueListenable<Shot> shot;
  final bool isSelected;
  final VoidCallback? onPressed;

  const RangeCardButton({
    super.key,
    required this.shot,
    required this.isSelected,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return ListenableBuilder(
      listenable: shot,
      builder: (context, _) {
        var button = isSelected ? FilledButton.new : FilledButton.tonal;
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 2.0),
          child: AspectRatio(
            aspectRatio: 1,
            child: button(
              onPressed: onPressed,
              child: Text(shot.value.toString(), textAlign: TextAlign.center),
              style: FilledButton.styleFrom(padding: EdgeInsets.all(0)),
            ),
          ),
        );
      },
    );
  }
}
